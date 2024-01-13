options=("deploy all")

select opt in "${options[@]}"
do
    case $opt in
        "deploy all")
            echo "input api version: "
            read version
            export API_VERSION=$version
            echo "input web version: "
            read version
            export WEB_VERSION=$version
            docker compose -f docker-compose.yaml up -d
            ;;
        *) echo "invalid option $REPLY";;
    esac
    break
done