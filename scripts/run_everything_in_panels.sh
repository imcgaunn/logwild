#!/bin/zsh

setopt err_exit
setopt no_unset
setopt pipe_fail

SCRIPT_DIR="${0:A:h}"

function die() {
  echo "did i die?"
}

function cleanup() {
	printf "killing session created\n"
	tmux kill-session -t observe_demo
	printf "successfully cleaned up\n"
}

# create a new tmux session
tmux new-session -d -s observe_demo
tmux send-keys -t observe_demo:1.1 'just run-observe-backend' C-m
# split window horizontally
tmux split-window -t observe_demo:1.1 -h
tmux send-keys -t observe_demo:1.2 'just run-otel-collector' C-m
# split window vertically
tmux split-window -t observe_demo:1.1 -v
tmux send-keys -t observe_demo:1.3 'just run-app' C-m
if [ -z "$TMUX" ]; then
  echo "tmux not set"
  tmux attach-session -t observe_demo
else
  echo "tmux set"
  tmux switchc -t observe_demo
fi

echo "do we know at this point that we have exited the new tmux sess?"
