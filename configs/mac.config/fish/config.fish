if status is-interactive
    # Commands to run in interactive sessions can go here
    set -U fish_user_paths $fish_user_paths /opt/homebrew/bin
    set -xU EDITOR nvim 
    alias cat='bat'
    alias ls='eza' 
    alias ll='eza -l'
    alias cs='cd'
    alias vim='nvim'
    alias run-cigarland='podman run -d --rm -v ./runtime:/cigarland -p 80:80 cigarland:latest'
    function build-cigarland;
      podman run --rm -v ~/git-projects/cigarland/go_app/:/go/src -w /go/src golang:1.24 /bin/sh -c "export CGO_ENABLED=1 && go build"
        #      export GOOS=linux
        #      export GOARCH=amd64 
        #      export CGO_ENABLED=1
        #      set cur_dir (pwd)
        #      cd ~/git-projects/cigarland/go_app 
        #      go build 
      ln  -f ~/git-projects/cigarland/go_app/cigarland_api ~/git-projects/cigarland/runtime/cigarland_api 
        #      set -e GOOS
        #      set -e GOARCH
        #      set -e CGO_ENABLED
        #      cd $cur_dir
    end
    function gocommit;
        if test (count $argv) = 0
            git commit *.go -m 'committing more work'
            git push origin HEAD:dev
        else if test (count $argv) = 1
            git commit *.go -m 'committing more work'
            git push origin HEAD:$argv[1]
        else if test (count $argv) > 1
            git commit *.go -m (string join ' ' $argv[2..-1])
            git push origin HEAD:$argv[1]
        else 
            git commit *.go -m 'committing more work'
            git push origin HEAD:$argv[1]
        end
    end
    tmux

end



