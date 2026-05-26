# App but better...
alias cl=clear
alias cat=bat
alias ll="eza -l --icons --git -a"
alias lt="eza --tree --level=2 --long --icons --git"
alias vim=nvim

# System update
alias update="sudo dnf update -y && flatpak update -y"

# Podman
alias pp="podman ps"
alias pcu="podman compose up"
alias pcd="podman compose down"
alias pl="podman logs"
alias pv="podman volume"
alias pe="podman exec"
alais pei="podman exec -it"
