if status is-interactive
    # Commands to run in interactive sessions can go here
    set -U fish_user_paths $fish_user_paths /opt/homebrew/bin
    alias cat='bat'
    alias ls='eza' 
    alias ll='eza -l'
    alias fze='find . -type f | fzf --bind "enter:execute(nvim {})+abort"'
    function copy-configs 
        set oldDir (pwd)
        cp -r ~/.config ~/git-projects/programming/configs/mac.config/ 
         cd ~/git-projects/programming/configs/mac.config/
         git add . 
         git commit -m 'updating configs' 
         git push origin HEAD:config-update
        cd $oldDir
    end
    tmux

end



