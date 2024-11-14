set -euo pipefail

baseWorkspace=${1:-'{}'}

echo "Base workspace file content:"
echo $baseWorkspace

# Get the folders from the base workspace file
folders=$(jq '.folders // []' <<<"$baseWorkspace")

# For each directory in /mnt/dev, add it to the workspace file
for dir in /mnt/dev/*; do
  if [ -d $dir ]; then
    # Get the base name of the directory
    basename=$(basename $dir)

    echo "Adding $basename to workspace file"

    # Add the directory to the folders array
    folders=$(jq --arg name "$basename" --arg path "$dir" '. + [{name: $name, path: $path}]' <<<"$folders")
  fi
done

# Create a new JSON object with the folders array
newJson=$(jq -n --argjson folders "$folders" '{folders: $folders}')

echo "Merging with base workspace"
# Merge the new JSON object with the base workspace file, that will fully replace the folders array, and keep the rest of the base workspace file
newJson=$(echo "$baseWorkspace $newJson" | jq -s add)

echo "New workspace file content:"
echo $newJson

echo "$newJson" >~/workspace.code-workspace
