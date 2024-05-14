#!/bin/zsh

setopt err_exit
setopt no_unset
setopt pipe_fail

declare -g LOCAL_BIN_DIR="$HOME/.local/bin"

function get() {
    local version="$1"
    # TODO: support other platforms here, I just don't want to bother
    local zinc_release_basename=$(printf "zincsearch_%s_linux_x86_64.tar.gz" "${version}")
    local zinc_releases_fmt="https://github.com/zincsearch/zincsearch/releases/download/v%s/%s"
    local zinc_release_url=$(printf "${zinc_releases_fmt}" "${version}" "${zinc_release_basename}")
    local zincsearch_exe="zincsearch"
    local zinc_local_bin_dest="$HOME/.local/bin/${zincsearch_exe}"
    [[ -d "${LOCAL_BIN_DIR}" ]] || die "can't install in ${LOCAL_BIN_DIR}; it doesn't exist"
    if [[ -x "${zinc_local_bin_dest}" ]]; then
        printf "zincsearch already installed at %s\n" "${zinc_local_bin_dest}"
        printf "nothing to be done!\n"
        return
    fi

    printf "downloading %s\n" "${zinc_release_url}"
    pushd /tmp
    curl -LO "${zinc_release_url}"
    tar xf "${zinc_release_basename}" "${zincsearch_exe}"
    mv "${zincsearch_exe}" "${zinc_local_bin_dest}"
}

get "0.4.10"
