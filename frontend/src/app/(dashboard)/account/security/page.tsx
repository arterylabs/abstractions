"use client";

import api from "@/api";
import useAPI from "@/api/hooks";
import { Button } from "@/components/ui/button";
import { SectionMark } from "@/components/ui/section-mark";
import { Input } from "@/components/ui/input";
import { Notifier } from "@/components/ui/sonner";
import { useAppSelector } from "@/store/hooks";
import { FormEvent, useState, useEffect } from "react";
import * as yup from "yup";
import Image from "next/image";
import { SimpleDialog } from "@/components/ui/dialog";
import Icon from "@/components/ui/icon";
import { Card, CardContent } from "@/components/ui/card";

const passwordSchema = yup.object().shape({
  newPassword: yup.string().min(8).required(),
  confirmPassword: yup
    .string()
    .oneOf([yup.ref("newPassword")], "Passwords must match")
    .required(),
});

export default function SecuritySettings() {
  const [currentPassword, setCurrentPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [isGoogleLinking, setIsGoogleLinking] = useState(false);
  const [isGitHubLinking, setIsGitHubLinking] = useState(false);

  const currentUser = useAppSelector((state) => state.auth.currentUser);

  const [showPasswordUpdateModal, setShowPasswordUpdateModal] = useState(false);

  const [updatePasswordResponse, updatePassword] = useAPI(() =>
    api.users.updatePassword({
      currentPassword,
      newPassword,
    }),
  );

  const handlePasswordUpdate = async (e: FormEvent) => {
    e.preventDefault();
    try {
      await passwordSchema.validate({
        currentPassword,
        newPassword,
        confirmPassword,
      });
      await updatePassword();
    } catch (error: unknown) {
      Notifier.error({
        title: "Error",
        message: (error as Error).message,
      });
    }
  };

  useEffect(() => {
    if (
      !updatePasswordResponse.error &&
      updatePasswordResponse.status === 200
    ) {
      setCurrentPassword("");
      setNewPassword("");
      setConfirmPassword("");
      setShowPasswordUpdateModal(false);
      Notifier.success({
        title: "Password updated",
        message: "Your password has been updated successfully",
      });

      // reload the user
      api.utils.removeToken();
      window.location.reload();
    }
    if (updatePasswordResponse.error) {
      Notifier.error({
        title: "Error",
        message: updatePasswordResponse.error.message,
      });
    }
  }, [updatePasswordResponse]);

  const handleGoogleLink = () => {
    setIsGoogleLinking(true);
    window.location.href = `${process.env.NEXT_PUBLIC_API_URL}/auth/google?ref=/account/security`;
  };

  const handleGitHubLink = () => {
    setIsGitHubLinking(true);
    window.location.href = `${process.env.NEXT_PUBLIC_API_URL}/auth/github?ref=/account/security`;
  };

  if (!currentUser) return null;

  const googleAuthenticatorExists =
    currentUser.authenticators?.includes("google");
  const githubAuthenticatorExists =
    currentUser.authenticators?.includes("github");
  const passwordAuthenticatorExists =
    currentUser.authenticators?.includes("password");

  const openPasswordUpdateModal = () => {
    setShowPasswordUpdateModal(true);
  };

  return (
    <section className="flex flex-col gap-5">
      <SectionMark label="Sign-in methods" />
      <Card>
        <CardContent className="flex flex-wrap items-center justify-between gap-4 p-6">
          <div className="flex items-center gap-4">
            <Image
              src="/assets/img/google.svg"
              alt="Google"
              width={30}
              height={30}
            />
            <div>
              <h3 className="text-xl font-medium">Google account</h3>
              <p className="mt-1 text-sm leading-relaxed text-muted-foreground">
                Sign in faster with your Google account.
              </p>
            </div>
          </div>

          {!googleAuthenticatorExists ? (
            <Button
              onClick={handleGoogleLink}
              variant="outline"
              loading={isGoogleLinking}
            >
              Connect Google Account
            </Button>
          ) : (
            <span className="press-label flex items-center gap-1.5 text-primary">
              <Icon type="check_circle" className="!text-sm" /> Connected
            </span>
          )}
        </CardContent>
      </Card>

      <Card>
        <CardContent className="flex flex-wrap items-center justify-between gap-4 p-6">
          <div className="flex items-center gap-4">
            <Image
              src="/assets/img/github.svg"
              alt="GitHub"
              className="invert"
              width={30}
              height={30}
            />
            <div>
              <h3 className="text-xl font-medium">GitHub account</h3>
              <p className="mt-1 text-sm leading-relaxed text-muted-foreground">
                Sign in faster with your GitHub account.
              </p>
            </div>
          </div>

          {!githubAuthenticatorExists ? (
            <Button
              onClick={handleGitHubLink}
              variant="outline"
              loading={isGitHubLinking}
            >
              Connect GitHub Account
            </Button>
          ) : (
            <span className="press-label flex items-center gap-1.5 text-primary">
              <Icon type="check_circle" className="!text-sm" /> Connected
            </span>
          )}
        </CardContent>
      </Card>

      <Card>
        <CardContent className="flex flex-wrap items-center justify-between gap-4 p-6">
          <div className="flex items-center gap-4">
            <Image
              src="/assets/img/key.svg"
              alt="Password"
              width={30}
              height={30}
            />
            <div>
              <h3 className="text-xl font-medium">Password</h3>
              <p className="mt-1 text-sm leading-relaxed text-muted-foreground">
                {passwordAuthenticatorExists
                  ? "A password is set for this account."
                  : "No password set — you're using OAuth only."}
              </p>
            </div>
          </div>

          {!passwordAuthenticatorExists ? (
            <Button
              onClick={openPasswordUpdateModal}
              variant="outline"
              data-testid="setup-password-button"
            >
              Setup Password
            </Button>
          ) : (
            <Button
              onClick={openPasswordUpdateModal}
              variant="outline"
              data-testid="change-password-button"
            >
              Change Password
            </Button>
          )}
        </CardContent>
      </Card>

      {showPasswordUpdateModal && (
        <SimpleDialog
          open={showPasswordUpdateModal}
          onOpenChange={setShowPasswordUpdateModal}
          title="Change Password"
          description="Update your account password"
          footer={
            <div className="flex justify-end gap-4 grow">
              <Button
                className="w-full"
                onClick={() => setShowPasswordUpdateModal(false)}
                variant="outline"
                disabled={updatePasswordResponse.isFetching}
              >
                Cancel
              </Button>
              <Button
                type="submit"
                className="w-full"
                onClick={handlePasswordUpdate}
                loading={updatePasswordResponse.isFetching}
              >
                {passwordAuthenticatorExists
                  ? "Update Password"
                  : "Set Password"}
              </Button>
            </div>
          }
        >
          <form onSubmit={handlePasswordUpdate} className="space-y-4">
            {passwordAuthenticatorExists && (
              <Input
                id="currentPassword"
                type="password"
                label="Current Password"
                placeholder="Enter your current password"
                value={currentPassword}
                onChange={setCurrentPassword}
                toggleVisibility
              />
            )}
            <div>
              <Input
                id="newPassword"
                type="password"
                label="New Password"
                placeholder="Enter your new password"
                value={newPassword}
                onChange={(v) => setNewPassword(v)}
                toggleVisibility
              />
            </div>
            <div>
              <Input
                id="confirmPassword"
                type="password"
                label="Confirm New Password"
                placeholder="Confirm your new password"
                value={confirmPassword}
                onChange={(v) => setConfirmPassword(v)}
                toggleVisibility
              />
            </div>
          </form>
        </SimpleDialog>
      )}
    </section>
  );
}
