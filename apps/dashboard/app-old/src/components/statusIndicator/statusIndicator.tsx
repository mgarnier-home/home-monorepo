type StatusIndicatorProps = {
  color: string;
  className?: string;
};

function StatusIndicator(props: StatusIndicatorProps) {
  const { color, className } = props;

  const style: React.CSSProperties = !color.startsWith("bg-") ? { backgroundColor: color } : {};

  return (
    <div
      className={`w-3 h-3 rounded-full ${className || ""} ${color.startsWith("bg-") ? color : ""}`}
      style={style}
    ></div>
  );
}

export default StatusIndicator;
