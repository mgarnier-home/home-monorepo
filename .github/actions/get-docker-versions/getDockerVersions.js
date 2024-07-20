const main = async () => {
  console.log(process.env);
  // const apiUrl = core.getInput('api-url');
  // const imageAuthor = core.getInput('image-author');
  // const imageName = core.getInput('image-name');

  // const onlyLatest = core.getInput('only-latest') === 'true';

  // console.log(apiUrl, imageAuthor, imageName, onlyLatest);
};

main().catch((error) => {
  console.error(error);
  process.exit(1);
});
