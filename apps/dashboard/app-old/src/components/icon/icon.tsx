const hlCdn = "https://raw.githubusercontent.com/walkxcode/dashboard-icons/master/png/{icon}.png";

type IconSize = "sm" | "md" | "lg";

type IconProps = {
  icon: string;
  size: IconSize;
};

function Icon(props: IconProps) {
  const { icon, size } = props;

  const renderImgIcon = (src: string): JSX.Element => {
    let className = "min-w-1rem min-h-1rem";

    if (size === "sm") className += " max-w-1.5rem max-h-1.5rem";
    else if (size === "md") className += " max-w-2rem max-h-2rem";
    else if (size === "lg") className += " max-w-3rem max-h-3rem";

    return <img src={src} className={className} />;
  };

  const renderFaIcon = (icon: string): JSX.Element => {
    let className = "mx-5px";

    if (size === "sm") className += " text-1.5rem";
    else if (size === "md") className += " text-2rem";
    else if (size === "lg") className += " text-2.5rem";

    return <i className={`${icon} ${className}`} />;
  };

  let iconRender = null;

  if (icon.startsWith("hl-")) {
    const src = hlCdn.replace("{icon}", icon.replace("hl-", ""));

    iconRender = renderImgIcon(src);
  } else if (icon.startsWith("fas")) {
    iconRender = renderFaIcon(icon);
  } else {
    iconRender = renderImgIcon(icon);
  }

  return <div className="flex items-center">{iconRender}</div>;
}

export default Icon;
