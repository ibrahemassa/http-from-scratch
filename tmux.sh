#!/bin/bash

tmux new-session -d -s HttpGo

tmux rename-window -t HttpGo:0 'nvim'

tmux send-keys -t HttpGo:1 'nvim' C-m

tmux new-window -t HttpGo:2 

tmux new-window -t HttpGo:3

tmux attach -t HttpGo

