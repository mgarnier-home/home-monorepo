FROM mgarnier11/my-home:deps

RUN pnpm run --filter "./libs/**" build

