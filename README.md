# pig in a blanket

run PiG (https://pi-in-go.dev/) in a container

## extras

- container has mise + golang installed
- comes with `/plan` mode extension (read only tools) enabled by default

## setup

1. ensure [mise](https://mise.jdx.dev/) + docker is installed
1. build the container `mise run build`
1. add to `~/.zshrc`, do not forget to source it
   ```sh
    function z-pig-here() {
        PIG_HOME=some-dir/pig-home ## where to persist all the PiG state: auth, sessions, tools, etc
        echo "running PiG here: $(pwd -P)"
        docker run --interactive --tty --rm \
            --env PI_TELEMETRY=0 \
            --env PI_SKIP_VERSION_CHECK=1 \
            --env PLAN_MODE_ENABLED=true \
            --volume $PIG_HOME:/home/pig-user/.pig \
            --volume "$(pwd -P)":/workspace \
            pig-in-a-blanket \
                --model github-copilot/claude-sonnet-5.5 \
                --thinking medium
    }
   ```
1. run `z-pig-here` in some directory
1. ...
1. profit
