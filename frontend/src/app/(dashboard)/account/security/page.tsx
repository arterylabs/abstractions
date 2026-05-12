"use client";

import api from "@/api";
import useAPI from "@/api/hooks";
import { Button } from "@/components/ui/button";
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
    window.location.href = `${process.env.NEXT_PUBLIC_API_URL}/auth/google/link`;
  };

  if (!currentUser) return null;

  const googleAuthenticatorExists =
    currentUser.authenticators?.includes("google");
  const passwordAuthenticatorExists =
    currentUser.authenticators?.includes("password");

  const openPasswordUpdateModal = () => {
    setShowPasswordUpdateModal(true);
  };

  return (
    <div className="space-y-4">
      <Card>
        <CardContent className="flex items-center justify-between p-6">
          <div className="flex items-center gap-4">
            <Image
              src="/assets/img/google.svg"
              alt="Google"
              width={30}
              height={30}
            />
            <div>
              <h3 className="font-medium">Google Account</h3>
              <p className="text-sm text-muted-foreground mt-1">
                Using Gmail for <b>{currentUser.email}</b>? You can login with
                Google too
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
            <div className="text-sm text-primary flex items-center justify-center gap-1">
              <Icon type="check_circle" /> Connected
            </div>
          )}
        </CardContent>
      </Card>

      <Card>
        <CardContent className="flex items-center justify-between p-6">
          <div className="flex items-center gap-4">
            <Image
              src="/assets/img/key.svg"
              alt="Password"
              width={30}
              height={30}
            />
            <div>
              <h3 className="font-medium">Password</h3>
              <p className="text-sm text-muted-foreground">
                {passwordAuthenticatorExists
                  ? "You have a password set for this account"
                  : "No password set for this account"}
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
    </div>
  );
}
