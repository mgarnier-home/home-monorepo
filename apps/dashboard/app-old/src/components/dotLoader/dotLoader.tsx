import "./dotLoader.css";

type DotLoaderProps = {
  nbDots: number;
  color: string;
  width?: string;
  height?: string;
  margin?: string;
};

function BouncingDotsLoader(props: DotLoaderProps) {
  const { nbDots, color, width, height, margin } = props;

  const style: React.CSSProperties = {
    backgroundColor: !color.startsWith("bg-") ? color : "",
    width: width || "0.4rem",
    height: height || "0.4rem",
    marginLeft: margin || (width ? `calc(${width} / 3)` : "0.08rem"),
    // marginRight: margin || (width ? `calc(${width} / 5)` : "0.08rem"),
  };

  return (
    <div className={`bouncing-loader`}>
      {[...Array(nbDots)].map((_, index) => (
        <div key={index} className={`dot ${color.startsWith("bg-") ? color : ""}`} style={style}></div>
      ))}
    </div>
  );
}

export default BouncingDotsLoader;

export const loaderSm = (color: string) => (
  <BouncingDotsLoader nbDots={3} color={color} width="0.19rem" height="0.19rem" />
);
