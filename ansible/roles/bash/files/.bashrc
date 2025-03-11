# .bashrc

HISTTIMEFORMAT='dd-mm-yyyy %T '
# Don't put duplicate lines in the history. See bash(1) for more options
export HISTCONTROL=ignoredups

# Source files
for file in $(find $HOME/.config/bash/ -type f -name "*.sh"); do
  source "$file"
done

# Starship
eval "$(starship init bash)"

