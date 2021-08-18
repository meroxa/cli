# bash completion for meroxa                               -*- shell-script -*-

__meroxa_debug()
{
    if [[ -n ${BASH_COMP_DEBUG_FILE} ]]; then
        echo "$*" >> "${BASH_COMP_DEBUG_FILE}"
    fi
}

# Homebrew on Macs have version 1.3 of bash-completion which doesn't include
# _init_completion. This is a very minimal version of that function.
__meroxa_init_completion()
{
    COMPREPLY=()
    _get_comp_words_by_ref "$@" cur prev words cword
}

__meroxa_index_of_word()
{
    local w word=$1
    shift
    index=0
    for w in "$@"; do
        [[ $w = "$word" ]] && return
        index=$((index+1))
    done
    index=-1
}

__meroxa_contains_word()
{
    local w word=$1; shift
    for w in "$@"; do
        [[ $w = "$word" ]] && return
    done
    return 1
}

__meroxa_handle_go_custom_completion()
{
    __meroxa_debug "${FUNCNAME[0]}: cur is ${cur}, words[*] is ${words[*]}, #words[@] is ${#words[@]}"

    local shellCompDirectiveError=1
    local shellCompDirectiveNoSpace=2
    local shellCompDirectiveNoFileComp=4
    local shellCompDirectiveFilterFileExt=8
    local shellCompDirectiveFilterDirs=16

    local out requestComp lastParam lastChar comp directive args

    # Prepare the command to request completions for the program.
    # Calling ${words[0]} instead of directly meroxa allows to handle aliases
    args=("${words[@]:1}")
    requestComp="${words[0]} __completeNoDesc ${args[*]}"

    lastParam=${words[$((${#words[@]}-1))]}
    lastChar=${lastParam:$((${#lastParam}-1)):1}
    __meroxa_debug "${FUNCNAME[0]}: lastParam ${lastParam}, lastChar ${lastChar}"

    if [ -z "${cur}" ] && [ "${lastChar}" != "=" ]; then
        # If the last parameter is complete (there is a space following it)
        # We add an extra empty parameter so we can indicate this to the go method.
        __meroxa_debug "${FUNCNAME[0]}: Adding extra empty parameter"
        requestComp="${requestComp} \"\""
    fi

    __meroxa_debug "${FUNCNAME[0]}: calling ${requestComp}"
    # Use eval to handle any environment variables and such
    out=$(eval "${requestComp}" 2>/dev/null)

    # Extract the directive integer at the very end of the output following a colon (:)
    directive=${out##*:}
    # Remove the directive
    out=${out%:*}
    if [ "${directive}" = "${out}" ]; then
        # There is not directive specified
        directive=0
    fi
    __meroxa_debug "${FUNCNAME[0]}: the completion directive is: ${directive}"
    __meroxa_debug "${FUNCNAME[0]}: the completions are: ${out[*]}"

    if [ $((directive & shellCompDirectiveError)) -ne 0 ]; then
        # Error code.  No completion.
        __meroxa_debug "${FUNCNAME[0]}: received error from custom completion go code"
        return
    else
        if [ $((directive & shellCompDirectiveNoSpace)) -ne 0 ]; then
            if [[ $(type -t compopt) = "builtin" ]]; then
                __meroxa_debug "${FUNCNAME[0]}: activating no space"
                compopt -o nospace
            fi
        fi
        if [ $((directive & shellCompDirectiveNoFileComp)) -ne 0 ]; then
            if [[ $(type -t compopt) = "builtin" ]]; then
                __meroxa_debug "${FUNCNAME[0]}: activating no file completion"
                compopt +o default
            fi
        fi
    fi

    if [ $((directive & shellCompDirectiveFilterFileExt)) -ne 0 ]; then
        # File extension filtering
        local fullFilter filter filteringCmd
        # Do not use quotes around the $out variable or else newline
        # characters will be kept.
        for filter in ${out[*]}; do
            fullFilter+="$filter|"
        done

        filteringCmd="_filedir $fullFilter"
        __meroxa_debug "File filtering command: $filteringCmd"
        $filteringCmd
    elif [ $((directive & shellCompDirectiveFilterDirs)) -ne 0 ]; then
        # File completion for directories only
        local subDir
        # Use printf to strip any trailing newline
        subdir=$(printf "%s" "${out[0]}")
        if [ -n "$subdir" ]; then
            __meroxa_debug "Listing directories in $subdir"
            __meroxa_handle_subdirs_in_dir_flag "$subdir"
        else
            __meroxa_debug "Listing directories in ."
            _filedir -d
        fi
    else
        while IFS='' read -r comp; do
            COMPREPLY+=("$comp")
        done < <(compgen -W "${out[*]}" -- "$cur")
    fi
}

__meroxa_handle_reply()
{
    __meroxa_debug "${FUNCNAME[0]}"
    local comp
    case $cur in
        -*)
            if [[ $(type -t compopt) = "builtin" ]]; then
                compopt -o nospace
            fi
            local allflags
            if [ ${#must_have_one_flag[@]} -ne 0 ]; then
                allflags=("${must_have_one_flag[@]}")
            else
                allflags=("${flags[*]} ${two_word_flags[*]}")
            fi
            while IFS='' read -r comp; do
                COMPREPLY+=("$comp")
            done < <(compgen -W "${allflags[*]}" -- "$cur")
            if [[ $(type -t compopt) = "builtin" ]]; then
                [[ "${COMPREPLY[0]}" == *= ]] || compopt +o nospace
            fi

            # complete after --flag=abc
            if [[ $cur == *=* ]]; then
                if [[ $(type -t compopt) = "builtin" ]]; then
                    compopt +o nospace
                fi

                local index flag
                flag="${cur%=*}"
                __meroxa_index_of_word "${flag}" "${flags_with_completion[@]}"
                COMPREPLY=()
                if [[ ${index} -ge 0 ]]; then
                    PREFIX=""
                    cur="${cur#*=}"
                    ${flags_completion[${index}]}
                    if [ -n "${ZSH_VERSION}" ]; then
                        # zsh completion needs --flag= prefix
                        eval "COMPREPLY=( \"\${COMPREPLY[@]/#/${flag}=}\" )"
                    fi
                fi
            fi
            return 0;
            ;;
    esac

    # check if we are handling a flag with special work handling
    local index
    __meroxa_index_of_word "${prev}" "${flags_with_completion[@]}"
    if [[ ${index} -ge 0 ]]; then
        ${flags_completion[${index}]}
        return
    fi

    # we are parsing a flag and don't have a special handler, no completion
    if [[ ${cur} != "${words[cword]}" ]]; then
        return
    fi

    local completions
    completions=("${commands[@]}")
    if [[ ${#must_have_one_noun[@]} -ne 0 ]]; then
        completions+=("${must_have_one_noun[@]}")
    elif [[ -n "${has_completion_function}" ]]; then
        # if a go completion function is provided, defer to that function
        __meroxa_handle_go_custom_completion
    fi
    if [[ ${#must_have_one_flag[@]} -ne 0 ]]; then
        completions+=("${must_have_one_flag[@]}")
    fi
    while IFS='' read -r comp; do
        COMPREPLY+=("$comp")
    done < <(compgen -W "${completions[*]}" -- "$cur")

    if [[ ${#COMPREPLY[@]} -eq 0 && ${#noun_aliases[@]} -gt 0 && ${#must_have_one_noun[@]} -ne 0 ]]; then
        while IFS='' read -r comp; do
            COMPREPLY+=("$comp")
        done < <(compgen -W "${noun_aliases[*]}" -- "$cur")
    fi

    if [[ ${#COMPREPLY[@]} -eq 0 ]]; then
		if declare -F __meroxa_custom_func >/dev/null; then
			# try command name qualified custom func
			__meroxa_custom_func
		else
			# otherwise fall back to unqualified for compatibility
			declare -F __custom_func >/dev/null && __custom_func
		fi
    fi

    # available in bash-completion >= 2, not always present on macOS
    if declare -F __ltrim_colon_completions >/dev/null; then
        __ltrim_colon_completions "$cur"
    fi

    # If there is only 1 completion and it is a flag with an = it will be completed
    # but we don't want a space after the =
    if [[ "${#COMPREPLY[@]}" -eq "1" ]] && [[ $(type -t compopt) = "builtin" ]] && [[ "${COMPREPLY[0]}" == --*= ]]; then
       compopt -o nospace
    fi
}

# The arguments should be in the form "ext1|ext2|extn"
__meroxa_handle_filename_extension_flag()
{
    local ext="$1"
    _filedir "@(${ext})"
}

__meroxa_handle_subdirs_in_dir_flag()
{
    local dir="$1"
    pushd "${dir}" >/dev/null 2>&1 && _filedir -d && popd >/dev/null 2>&1 || return
}

__meroxa_handle_flag()
{
    __meroxa_debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"

    # if a command required a flag, and we found it, unset must_have_one_flag()
    local flagname=${words[c]}
    local flagvalue
    # if the word contained an =
    if [[ ${words[c]} == *"="* ]]; then
        flagvalue=${flagname#*=} # take in as flagvalue after the =
        flagname=${flagname%=*} # strip everything after the =
        flagname="${flagname}=" # but put the = back
    fi
    __meroxa_debug "${FUNCNAME[0]}: looking for ${flagname}"
    if __meroxa_contains_word "${flagname}" "${must_have_one_flag[@]}"; then
        must_have_one_flag=()
    fi

    # if you set a flag which only applies to this command, don't show subcommands
    if __meroxa_contains_word "${flagname}" "${local_nonpersistent_flags[@]}"; then
      commands=()
    fi

    # keep flag value with flagname as flaghash
    # flaghash variable is an associative array which is only supported in bash > 3.
    if [[ -z "${BASH_VERSION}" || "${BASH_VERSINFO[0]}" -gt 3 ]]; then
        if [ -n "${flagvalue}" ] ; then
            flaghash[${flagname}]=${flagvalue}
        elif [ -n "${words[ $((c+1)) ]}" ] ; then
            flaghash[${flagname}]=${words[ $((c+1)) ]}
        else
            flaghash[${flagname}]="true" # pad "true" for bool flag
        fi
    fi

    # skip the argument to a two word flag
    if [[ ${words[c]} != *"="* ]] && __meroxa_contains_word "${words[c]}" "${two_word_flags[@]}"; then
			  __meroxa_debug "${FUNCNAME[0]}: found a flag ${words[c]}, skip the next argument"
        c=$((c+1))
        # if we are looking for a flags value, don't show commands
        if [[ $c -eq $cword ]]; then
            commands=()
        fi
    fi

    c=$((c+1))

}

__meroxa_handle_noun()
{
    __meroxa_debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"

    if __meroxa_contains_word "${words[c]}" "${must_have_one_noun[@]}"; then
        must_have_one_noun=()
    elif __meroxa_contains_word "${words[c]}" "${noun_aliases[@]}"; then
        must_have_one_noun=()
    fi

    nouns+=("${words[c]}")
    c=$((c+1))
}

__meroxa_handle_command()
{
    __meroxa_debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"

    local next_command
    if [[ -n ${last_command} ]]; then
        next_command="_${last_command}_${words[c]//:/__}"
    else
        if [[ $c -eq 0 ]]; then
            next_command="_meroxa_root_command"
        else
            next_command="_${words[c]//:/__}"
        fi
    fi
    c=$((c+1))
    __meroxa_debug "${FUNCNAME[0]}: looking for ${next_command}"
    declare -F "$next_command" >/dev/null && $next_command
}

__meroxa_handle_word()
{
    if [[ $c -ge $cword ]]; then
        __meroxa_handle_reply
        return
    fi
    __meroxa_debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"
    if [[ "${words[c]}" == -* ]]; then
        __meroxa_handle_flag
    elif __meroxa_contains_word "${words[c]}" "${commands[@]}"; then
        __meroxa_handle_command
    elif [[ $c -eq 0 ]]; then
        __meroxa_handle_command
    elif __meroxa_contains_word "${words[c]}" "${command_aliases[@]}"; then
        # aliashash variable is an associative array which is only supported in bash > 3.
        if [[ -z "${BASH_VERSION}" || "${BASH_VERSINFO[0]}" -gt 3 ]]; then
            words[c]=${aliashash[${words[c]}]}
            __meroxa_handle_command
        else
            __meroxa_handle_noun
        fi
    else
        __meroxa_handle_noun
    fi
    __meroxa_handle_word
}

_meroxa_api()
{
    last_command="meroxa_api"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--help")
    flags+=("-h")
    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--debug")
    flags+=("--json")
    flags+=("--timeout=")
    two_word_flags+=("--timeout")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_meroxa_auth_help()
{
    last_command="meroxa_auth_help"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--debug")
    flags+=("--json")
    flags+=("--timeout=")
    two_word_flags+=("--timeout")

    must_have_one_flag=()
    must_have_one_noun=()
    has_completion_function=1
    noun_aliases=()
}

_meroxa_auth_login()
{
    last_command="meroxa_auth_login"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--help")
    flags+=("-h")
    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--debug")
    flags+=("--json")
    flags+=("--timeout=")
    two_word_flags+=("--timeout")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_meroxa_auth_logout()
{
    last_command="meroxa_auth_logout"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--help")
    flags+=("-h")
    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--debug")
    flags+=("--json")
    flags+=("--timeout=")
    two_word_flags+=("--timeout")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_meroxa_auth_whoami()
{
    last_command="meroxa_auth_whoami"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--help")
    flags+=("-h")
    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--debug")
    flags+=("--json")
    flags+=("--timeout=")
    two_word_flags+=("--timeout")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_meroxa_auth()
{
    last_command="meroxa_auth"

    command_aliases=()

    commands=()
    commands+=("help")
    commands+=("login")
    commands+=("logout")
    commands+=("whoami")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--help")
    flags+=("-h")
    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--debug")
    flags+=("--json")
    flags+=("--timeout=")
    two_word_flags+=("--timeout")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_meroxa_billing()
{
    last_command="meroxa_billing"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--help")
    flags+=("-h")
    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--debug")
    flags+=("--json")
    flags+=("--timeout=")
    two_word_flags+=("--timeout")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_meroxa_completion()
{
    last_command="meroxa_completion"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--help")
    flags+=("-h")
    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--debug")
    flags+=("--json")
    flags+=("--timeout=")
    two_word_flags+=("--timeout")

    must_have_one_flag=()
    must_have_one_noun=()
    must_have_one_noun+=("bash")
    must_have_one_noun+=("fish")
    must_have_one_noun+=("powershell")
    must_have_one_noun+=("zsh")
    noun_aliases=()
}

_meroxa_config_describe()
{
    last_command="meroxa_config_describe"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--help")
    flags+=("-h")
    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--debug")
    flags+=("--json")
    flags+=("--timeout=")
    two_word_flags+=("--timeout")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_meroxa_config_help()
{
    last_command="meroxa_config_help"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--debug")
    flags+=("--json")
    flags+=("--timeout=")
    two_word_flags+=("--timeout")

    must_have_one_flag=()
    must_have_one_noun=()
    has_completion_function=1
    noun_aliases=()
}

_meroxa_config()
{
    last_command="meroxa_config"

    command_aliases=()

    commands=()
    commands+=("describe")
    commands+=("help")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--help")
    flags+=("-h")
    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--debug")
    flags+=("--json")
    flags+=("--timeout=")
    two_word_flags+=("--timeout")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_meroxa_connect()
{
    last_command="meroxa_connect"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--from=")
    two_word_flags+=("--from")
    flags+=("--help")
    flags+=("-h")
    flags+=("--input=")
    two_word_flags+=("--input")
    flags+=("--pipeline=")
    two_word_flags+=("--pipeline")
    flags+=("--to=")
    two_word_flags+=("--to")
    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--debug")
    flags+=("--json")
    flags+=("--timeout=")
    two_word_flags+=("--timeout")

    must_have_one_flag=()
    must_have_one_flag+=("--from=")
    must_have_one_flag+=("--to=")
    must_have_one_noun=()
    noun_aliases=()
}

_meroxa_connectors_create()
{
    last_command="meroxa_connectors_create"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--from=")
    two_word_flags+=("--from")
    flags+=("--help")
    flags+=("-h")
    flags+=("--input=")
    two_word_flags+=("--input")
    flags+=("--pipeline=")
    two_word_flags+=("--pipeline")
    flags+=("--to=")
    two_word_flags+=("--to")
    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--debug")
    flags+=("--json")
    flags+=("--timeout=")
    two_word_flags+=("--timeout")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_meroxa_connectors_describe()
{
    last_command="meroxa_connectors_describe"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--help")
    flags+=("-h")
    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--debug")
    flags+=("--json")
    flags+=("--timeout=")
    two_word_flags+=("--timeout")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_meroxa_connectors_help()
{
    last_command="meroxa_connectors_help"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--debug")
    flags+=("--json")
    flags+=("--timeout=")
    two_word_flags+=("--timeout")

    must_have_one_flag=()
    must_have_one_noun=()
    has_completion_function=1
    noun_aliases=()
}

_meroxa_connectors_list()
{
    last_command="meroxa_connectors_list"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--help")
    flags+=("-h")
    flags+=("--pipeline=")
    two_word_flags+=("--pipeline")
    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--debug")
    flags+=("--json")
    flags+=("--timeout=")
    two_word_flags+=("--timeout")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_meroxa_connectors_logs()
{
    last_command="meroxa_connectors_logs"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--help")
    flags+=("-h")
    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--debug")
    flags+=("--json")
    flags+=("--timeout=")
    two_word_flags+=("--timeout")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_meroxa_connectors_remove()
{
    last_command="meroxa_connectors_remove"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--force")
    flags+=("-f")
    flags+=("--help")
    flags+=("-h")
    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--debug")
    flags+=("--json")
    flags+=("--timeout=")
    two_word_flags+=("--timeout")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_meroxa_connectors_update()
{
    last_command="meroxa_connectors_update"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--help")
    flags+=("-h")
    flags+=("--state=")
    two_word_flags+=("--state")
    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--debug")
    flags+=("--json")
    flags+=("--timeout=")
    two_word_flags+=("--timeout")

    must_have_one_flag=()
    must_have_one_flag+=("--state=")
    must_have_one_noun=()
    noun_aliases=()
}

_meroxa_connectors()
{
    last_command="meroxa_connectors"

    command_aliases=()

    commands=()
    commands+=("create")
    commands+=("describe")
    commands+=("help")
    commands+=("list")
    if [[ -z "${BASH_VERSION}" || "${BASH_VERSINFO[0]}" -gt 3 ]]; then
        command_aliases+=("ls")
        aliashash["ls"]="list"
    fi
    commands+=("logs")
    commands+=("remove")
    if [[ -z "${BASH_VERSION}" || "${BASH_VERSINFO[0]}" -gt 3 ]]; then
        command_aliases+=("delete")
        aliashash["delete"]="remove"
        command_aliases+=("rm")
        aliashash["rm"]="remove"
    fi
    commands+=("update")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--help")
    flags+=("-h")
    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--debug")
    flags+=("--json")
    flags+=("--timeout=")
    two_word_flags+=("--timeout")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_meroxa_endpoints_create()
{
    last_command="meroxa_endpoints_create"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--help")
    flags+=("-h")
    flags+=("--protocol=")
    two_word_flags+=("--protocol")
    two_word_flags+=("-p")
    flags+=("--stream=")
    two_word_flags+=("--stream")
    two_word_flags+=("-s")
    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--debug")
    flags+=("--json")
    flags+=("--timeout=")
    two_word_flags+=("--timeout")

    must_have_one_flag=()
    must_have_one_flag+=("--protocol=")
    must_have_one_flag+=("-p")
    must_have_one_flag+=("--stream=")
    must_have_one_flag+=("-s")
    must_have_one_noun=()
    noun_aliases=()
}

_meroxa_endpoints_describe()
{
    last_command="meroxa_endpoints_describe"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--help")
    flags+=("-h")
    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--debug")
    flags+=("--json")
    flags+=("--timeout=")
    two_word_flags+=("--timeout")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_meroxa_endpoints_help()
{
    last_command="meroxa_endpoints_help"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--debug")
    flags+=("--json")
    flags+=("--timeout=")
    two_word_flags+=("--timeout")

    must_have_one_flag=()
    must_have_one_noun=()
    has_completion_function=1
    noun_aliases=()
}

_meroxa_endpoints_list()
{
    last_command="meroxa_endpoints_list"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--help")
    flags+=("-h")
    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--debug")
    flags+=("--json")
    flags+=("--timeout=")
    two_word_flags+=("--timeout")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_meroxa_endpoints_remove()
{
    last_command="meroxa_endpoints_remove"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--force")
    flags+=("-f")
    flags+=("--help")
    flags+=("-h")
    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--debug")
    flags+=("--json")
    flags+=("--timeout=")
    two_word_flags+=("--timeout")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_meroxa_endpoints()
{
    last_command="meroxa_endpoints"

    command_aliases=()

    commands=()
    commands+=("create")
    commands+=("describe")
    commands+=("help")
    commands+=("list")
    if [[ -z "${BASH_VERSION}" || "${BASH_VERSINFO[0]}" -gt 3 ]]; then
        command_aliases+=("ls")
        aliashash["ls"]="list"
    fi
    commands+=("remove")
    if [[ -z "${BASH_VERSION}" || "${BASH_VERSINFO[0]}" -gt 3 ]]; then
        command_aliases+=("delete")
        aliashash["delete"]="remove"
        command_aliases+=("rm")
        aliashash["rm"]="remove"
    fi

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--help")
    flags+=("-h")
    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--debug")
    flags+=("--json")
    flags+=("--timeout=")
    two_word_flags+=("--timeout")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_meroxa_help()
{
    last_command="meroxa_help"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--debug")
    flags+=("--json")
    flags+=("--timeout=")
    two_word_flags+=("--timeout")

    must_have_one_flag=()
    must_have_one_noun=()
    has_completion_function=1
    noun_aliases=()
}

_meroxa_login()
{
    last_command="meroxa_login"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--help")
    flags+=("-h")
    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--debug")
    flags+=("--json")
    flags+=("--timeout=")
    two_word_flags+=("--timeout")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_meroxa_logout()
{
    last_command="meroxa_logout"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--help")
    flags+=("-h")
    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--debug")
    flags+=("--json")
    flags+=("--timeout=")
    two_word_flags+=("--timeout")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_meroxa_open_billing()
{
    last_command="meroxa_open_billing"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--help")
    flags+=("-h")
    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--debug")
    flags+=("--json")
    flags+=("--timeout=")
    two_word_flags+=("--timeout")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_meroxa_open_help()
{
    last_command="meroxa_open_help"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--debug")
    flags+=("--json")
    flags+=("--timeout=")
    two_word_flags+=("--timeout")

    must_have_one_flag=()
    must_have_one_noun=()
    has_completion_function=1
    noun_aliases=()
}

_meroxa_open()
{
    last_command="meroxa_open"

    command_aliases=()

    commands=()
    commands+=("billing")
    commands+=("help")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--help")
    flags+=("-h")
    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--debug")
    flags+=("--json")
    flags+=("--timeout=")
    two_word_flags+=("--timeout")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_meroxa_pipelines_create()
{
    last_command="meroxa_pipelines_create"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--help")
    flags+=("-h")
    flags+=("--metadata=")
    two_word_flags+=("--metadata")
    two_word_flags+=("-m")
    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--debug")
    flags+=("--json")
    flags+=("--timeout=")
    two_word_flags+=("--timeout")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_meroxa_pipelines_help()
{
    last_command="meroxa_pipelines_help"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--debug")
    flags+=("--json")
    flags+=("--timeout=")
    two_word_flags+=("--timeout")

    must_have_one_flag=()
    must_have_one_noun=()
    has_completion_function=1
    noun_aliases=()
}

_meroxa_pipelines_list()
{
    last_command="meroxa_pipelines_list"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--help")
    flags+=("-h")
    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--debug")
    flags+=("--json")
    flags+=("--timeout=")
    two_word_flags+=("--timeout")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_meroxa_pipelines_remove()
{
    last_command="meroxa_pipelines_remove"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--force")
    flags+=("-f")
    flags+=("--help")
    flags+=("-h")
    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--debug")
    flags+=("--json")
    flags+=("--timeout=")
    two_word_flags+=("--timeout")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_meroxa_pipelines_update()
{
    last_command="meroxa_pipelines_update"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--help")
    flags+=("-h")
    flags+=("--metadata=")
    two_word_flags+=("--metadata")
    two_word_flags+=("-m")
    flags+=("--name=")
    two_word_flags+=("--name")
    flags+=("--state=")
    two_word_flags+=("--state")
    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--debug")
    flags+=("--json")
    flags+=("--timeout=")
    two_word_flags+=("--timeout")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_meroxa_pipelines()
{
    last_command="meroxa_pipelines"

    command_aliases=()

    commands=()
    commands+=("create")
    commands+=("help")
    commands+=("list")
    if [[ -z "${BASH_VERSION}" || "${BASH_VERSINFO[0]}" -gt 3 ]]; then
        command_aliases+=("ls")
        aliashash["ls"]="list"
    fi
    commands+=("remove")
    if [[ -z "${BASH_VERSION}" || "${BASH_VERSINFO[0]}" -gt 3 ]]; then
        command_aliases+=("delete")
        aliashash["delete"]="remove"
        command_aliases+=("rm")
        aliashash["rm"]="remove"
    fi
    commands+=("update")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--help")
    flags+=("-h")
    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--debug")
    flags+=("--json")
    flags+=("--timeout=")
    two_word_flags+=("--timeout")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_meroxa_resources_create()
{
    last_command="meroxa_resources_create"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--ca-cert=")
    two_word_flags+=("--ca-cert")
    flags+=("--client-cert=")
    two_word_flags+=("--client-cert")
    flags+=("--client-key=")
    two_word_flags+=("--client-key")
    flags+=("--help")
    flags+=("-h")
    flags+=("--metadata=")
    two_word_flags+=("--metadata")
    two_word_flags+=("-m")
    flags+=("--password=")
    two_word_flags+=("--password")
    flags+=("--ssh-private-key=")
    two_word_flags+=("--ssh-private-key")
    flags+=("--ssh-url=")
    two_word_flags+=("--ssh-url")
    flags+=("--ssl")
    flags+=("--type=")
    two_word_flags+=("--type")
    flags+=("--url=")
    two_word_flags+=("--url")
    two_word_flags+=("-u")
    flags+=("--username=")
    two_word_flags+=("--username")
    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--debug")
    flags+=("--json")
    flags+=("--timeout=")
    two_word_flags+=("--timeout")

    must_have_one_flag=()
    must_have_one_flag+=("--type=")
    must_have_one_flag+=("--url=")
    must_have_one_flag+=("-u")
    must_have_one_noun=()
    noun_aliases=()
}

_meroxa_resources_describe()
{
    last_command="meroxa_resources_describe"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--help")
    flags+=("-h")
    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--debug")
    flags+=("--json")
    flags+=("--timeout=")
    two_word_flags+=("--timeout")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_meroxa_resources_help()
{
    last_command="meroxa_resources_help"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--debug")
    flags+=("--json")
    flags+=("--timeout=")
    two_word_flags+=("--timeout")

    must_have_one_flag=()
    must_have_one_noun=()
    has_completion_function=1
    noun_aliases=()
}

_meroxa_resources_list()
{
    last_command="meroxa_resources_list"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--help")
    flags+=("-h")
    flags+=("--types")
    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--debug")
    flags+=("--json")
    flags+=("--timeout=")
    two_word_flags+=("--timeout")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_meroxa_resources_remove()
{
    last_command="meroxa_resources_remove"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--force")
    flags+=("-f")
    flags+=("--help")
    flags+=("-h")
    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--debug")
    flags+=("--json")
    flags+=("--timeout=")
    two_word_flags+=("--timeout")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_meroxa_resources_rotate-tunnel-key()
{
    last_command="meroxa_resources_rotate-tunnel-key"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--force")
    flags+=("-f")
    flags+=("--help")
    flags+=("-h")
    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--debug")
    flags+=("--json")
    flags+=("--timeout=")
    two_word_flags+=("--timeout")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_meroxa_resources_update()
{
    last_command="meroxa_resources_update"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--ca-cert=")
    two_word_flags+=("--ca-cert")
    flags+=("--client-cert=")
    two_word_flags+=("--client-cert")
    flags+=("--client-key=")
    two_word_flags+=("--client-key")
    flags+=("--help")
    flags+=("-h")
    flags+=("--metadata=")
    two_word_flags+=("--metadata")
    two_word_flags+=("-m")
    flags+=("--name=")
    two_word_flags+=("--name")
    flags+=("--password=")
    two_word_flags+=("--password")
    flags+=("--ssh-url=")
    two_word_flags+=("--ssh-url")
    flags+=("--ssl")
    flags+=("--url=")
    two_word_flags+=("--url")
    two_word_flags+=("-u")
    flags+=("--username=")
    two_word_flags+=("--username")
    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--debug")
    flags+=("--json")
    flags+=("--timeout=")
    two_word_flags+=("--timeout")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_meroxa_resources_validate()
{
    last_command="meroxa_resources_validate"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--help")
    flags+=("-h")
    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--debug")
    flags+=("--json")
    flags+=("--timeout=")
    two_word_flags+=("--timeout")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_meroxa_resources()
{
    last_command="meroxa_resources"

    command_aliases=()

    commands=()
    commands+=("create")
    if [[ -z "${BASH_VERSION}" || "${BASH_VERSINFO[0]}" -gt 3 ]]; then
        command_aliases+=("add")
        aliashash["add"]="create"
    fi
    commands+=("describe")
    commands+=("help")
    commands+=("list")
    if [[ -z "${BASH_VERSION}" || "${BASH_VERSINFO[0]}" -gt 3 ]]; then
        command_aliases+=("ls")
        aliashash["ls"]="list"
    fi
    commands+=("remove")
    if [[ -z "${BASH_VERSION}" || "${BASH_VERSINFO[0]}" -gt 3 ]]; then
        command_aliases+=("delete")
        aliashash["delete"]="remove"
        command_aliases+=("rm")
        aliashash["rm"]="remove"
    fi
    commands+=("rotate-tunnel-key")
    commands+=("update")
    commands+=("validate")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--help")
    flags+=("-h")
    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--debug")
    flags+=("--json")
    flags+=("--timeout=")
    two_word_flags+=("--timeout")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_meroxa_transforms_help()
{
    last_command="meroxa_transforms_help"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--debug")
    flags+=("--json")
    flags+=("--timeout=")
    two_word_flags+=("--timeout")

    must_have_one_flag=()
    must_have_one_noun=()
    has_completion_function=1
    noun_aliases=()
}

_meroxa_transforms_list()
{
    last_command="meroxa_transforms_list"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--help")
    flags+=("-h")
    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--debug")
    flags+=("--json")
    flags+=("--timeout=")
    two_word_flags+=("--timeout")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_meroxa_transforms()
{
    last_command="meroxa_transforms"

    command_aliases=()

    commands=()
    commands+=("help")
    commands+=("list")
    if [[ -z "${BASH_VERSION}" || "${BASH_VERSINFO[0]}" -gt 3 ]]; then
        command_aliases+=("ls")
        aliashash["ls"]="list"
    fi

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--help")
    flags+=("-h")
    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--debug")
    flags+=("--json")
    flags+=("--timeout=")
    two_word_flags+=("--timeout")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_meroxa_version()
{
    last_command="meroxa_version"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--help")
    flags+=("-h")
    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--debug")
    flags+=("--json")
    flags+=("--timeout=")
    two_word_flags+=("--timeout")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_meroxa_whoami()
{
    last_command="meroxa_whoami"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--help")
    flags+=("-h")
    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--debug")
    flags+=("--json")
    flags+=("--timeout=")
    two_word_flags+=("--timeout")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_meroxa_root_command()
{
    last_command="meroxa"

    command_aliases=()

    commands=()
    commands+=("api")
    commands+=("auth")
    commands+=("billing")
    commands+=("completion")
    commands+=("config")
    if [[ -z "${BASH_VERSION}" || "${BASH_VERSINFO[0]}" -gt 3 ]]; then
        command_aliases+=("cfg")
        aliashash["cfg"]="config"
    fi
    commands+=("connect")
    commands+=("connectors")
    if [[ -z "${BASH_VERSION}" || "${BASH_VERSINFO[0]}" -gt 3 ]]; then
        command_aliases+=("connector")
        aliashash["connector"]="connectors"
    fi
    commands+=("endpoints")
    if [[ -z "${BASH_VERSION}" || "${BASH_VERSINFO[0]}" -gt 3 ]]; then
        command_aliases+=("endpoint")
        aliashash["endpoint"]="endpoints"
    fi
    commands+=("help")
    commands+=("login")
    commands+=("logout")
    commands+=("open")
    commands+=("pipelines")
    if [[ -z "${BASH_VERSION}" || "${BASH_VERSINFO[0]}" -gt 3 ]]; then
        command_aliases+=("pipeline")
        aliashash["pipeline"]="pipelines"
    fi
    commands+=("resources")
    if [[ -z "${BASH_VERSION}" || "${BASH_VERSINFO[0]}" -gt 3 ]]; then
        command_aliases+=("resource")
        aliashash["resource"]="resources"
    fi
    commands+=("transforms")
    if [[ -z "${BASH_VERSION}" || "${BASH_VERSINFO[0]}" -gt 3 ]]; then
        command_aliases+=("transform")
        aliashash["transform"]="transforms"
    fi
    commands+=("version")
    commands+=("whoami")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--debug")
    flags+=("--help")
    flags+=("-h")
    flags+=("--json")
    flags+=("--timeout=")
    two_word_flags+=("--timeout")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

__start_meroxa()
{
    local cur prev words cword split
    declare -A flaghash 2>/dev/null || :
    declare -A aliashash 2>/dev/null || :
    if declare -F _init_completion >/dev/null 2>&1; then
        _init_completion -s || return
    else
        __meroxa_init_completion -n "=" || return
    fi

    local c=0
    local flags=()
    local two_word_flags=()
    local local_nonpersistent_flags=()
    local flags_with_completion=()
    local flags_completion=()
    local commands=("meroxa")
    local command_aliases=()
    local must_have_one_flag=()
    local must_have_one_noun=()
    local has_completion_function
    local last_command
    local nouns=()
    local noun_aliases=()

    __meroxa_handle_word
}

if [[ $(type -t compopt) = "builtin" ]]; then
    complete -o default -F __start_meroxa meroxa
else
    complete -o default -o nospace -F __start_meroxa meroxa
fi

# ex: ts=4 sw=4 et filetype=sh
