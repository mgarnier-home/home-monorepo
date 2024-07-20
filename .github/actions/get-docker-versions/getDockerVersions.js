const main = async () => {
  const apiUrl = process.env['INPUT_API-URL'];
  const imageAuthor = process.env['INPUT_IMAGE-AUTHOR'];
  const imageName = process.env['INPUT_IMAGE-NAME'];

  const onlyLatest = process.env['INPUT_ONLY-LATEST'] === 'true';

  console.log(apiUrl, imageAuthor, imageName, onlyLatest);
};

main().catch((error) => {
  console.error(error);
  process.exit(1);
});
