while read -r line; do
  ./latest.sh "https://beta.$line"
done < sites.txt