if status is-interactive
    # Commands to run in interactive sessions can go here
    set -U fish_user_paths $fish_user_paths /opt/homebrew/bin
    alias cat='bat'
    alias ls='eza' 
    alias ll='eza -l'
    tmux

end



