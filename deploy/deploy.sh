options=("deploy all" "shutdown all")

# get current api version
export API_VERSION=$(docker images | grep "api" | awk '{print $2}')
# get current web version
export WEB_VERSION=$(docker images | grep "web" | awk '{print $2}')

select opt in "${options[@]}"
do
    case $opt in
        "deploy all")
            echo "input api version: "
            read version
            API_VERSION=$version
            echo "input web version: "
            read version
            WEB_VERSION=$version
            docker compose -f docker-compose.yaml up -d
            ;;
        "shutdown all")
            docker compose -f docker-compose.yaml down chatgpt-plus-api chatgpt-plus-web
            ;;
        *) echo "invalid option $REPLY";;
    esac
    break
done