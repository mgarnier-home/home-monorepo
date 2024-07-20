const getAllTags = async (url) => {
  console.log(`Fetching ${url}`);

  const response = await fetch(url);
  const data = await response.json();

  if (data.next) {
    return [...data.results, ...(await getVersions(data.next))];
  } else {
    return data.results;
  }
};

const checkIfSemver = (tag) => {
  // comes from : https://semver.org
  const semverRegex =
    /^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$/;

  return semverRegex.test(tag);
};

const getVersions = async (url) => {
  const tags = await getAllTags(url);

  return tags.map((tag) => tag.name).filter((tag) => checkIfSemver(tag));
};

const main = async () => {
  const apiUrl = process.env['INPUT_API-URL'];
  const imageAuthor = process.env['INPUT_IMAGE-AUTHOR'];
  const imageName = process.env['INPUT_IMAGE-NAME'];

  const onlyLatest = process.env['INPUT_ONLY-LATEST'] === 'true';

  const url = `${apiUrl}/${imageAuthor}/${imageName}/tags`;

  const versions = await getVersions(url);

  console.log(versions);
};

main().catch((error) => {
  console.error(error);
  process.exit(1);
});
