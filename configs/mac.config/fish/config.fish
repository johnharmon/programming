if status is-interactive
    # Commands to run in interactive sessions can go here
    set -U fish_user_paths $fish_user_paths /opt/homebrew/bin
    set -xU EDITOR nvim 
    alias cat='bat'
    alias ls='eza' 
    alias ll='eza -l'
    alias cs='cd'
    alias vim='nvim'
    function gocommit;
        if test (count $argv) = 0
            git commit *.go -m 'committing more work'
            git push origin HEAD:dev
        else 
            git commit *.go -m 'committing more work'
            git push origin HEAD:$argv[1]
        end
    end
    tmux

end



