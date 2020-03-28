for filename in "$@"; do
    curl -X POST localhost:3000/images -H 'Content-Type: multipart/form-data' -F "file=@$(readlink -e $filename)"
done
