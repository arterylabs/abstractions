import fs from "fs";
import path from "path";
import { ImageResponse } from "next/og";

export const alt = "Abstractions — Keep your codebase knowledge fresh";
export const size = { width: 1200, height: 630 };
export const contentType = "image/png";

export default async function Image() {
  const logoData = fs.readFileSync(
    path.join(process.cwd(), "public/assets/logo.svg"),
  );
  const logoBase64 = `data:image/svg+xml;base64,${logoData.toString("base64")}`;

  return new ImageResponse(
    <div
      style={{
        background: "#09090b",
        width: "100%",
        height: "100%",
        display: "flex",
        flexDirection: "column",
        alignItems: "flex-start",
        justifyContent: "flex-end",
        padding: "64px",
      }}
    >
      {/* eslint-disable-next-line @next/next/no-img-element */}
      <img
        src={logoBase64}
        alt="Abstractions"
        style={{ height: 40, marginBottom: 48, objectFit: "contain" }}
      />
    </div>,
    { ...size },
  );
}
