"use client";

import api from "@/api";
import Spinner from "@/components/ui/spinner";
import { useEffect, useState } from "react";

export default function Authenticated() {
  const [token, setToken] = useState("");

  useEffect(() => {
    const token = api.utils.getToken();
    if (!token) {
      // If no token is found, redirect to login
      window.location.href = "/auth/login";
      return;
    }
    setToken(token);
  }, []);

  return (
    <div className="flex flex-col gap-4 items-center justify-center text-center">
      <Spinner />
      <div className="opacity-60 text-sm">
        Logging you in to the Abstractions application...
      </div>
      <div className="text-sm">
        <a href={`abstractions://?login=${token}`}>
          Click here if you are not automatically redirected
        </a>
      </div>
    </div>
  );
}
