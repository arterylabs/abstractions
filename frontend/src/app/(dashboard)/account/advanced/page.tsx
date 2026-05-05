"use client";

import { useState } from "react";
import api from "@/api";
import useAPI from "@/api/hooks";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Notifier } from "@/components/ui/sonner";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import { Card, CardContent } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import Icon from "@/components/ui/icon";
import { useRouter } from "next/navigation";
import { useAppSelector } from "@/store/hooks";

export default function AdvancedSettings() {
  const router = useRouter();
  const currentUser = useAppSelector((state) => state.auth.currentUser);
  const [, deleteAccount] = useAPI(() => api.users.deleteAccount());
  const [deleteConfirmText, setDeleteConfirmText] = useState("");
  const [isDialogOpen, setIsDialogOpen] = useState(false);

  const handleDeleteAccount = async () => {
    try {
      await deleteAccount();
      api.utils.removeToken();
      router.push("/login");
    } catch (error: unknown) {
      Notifier.error({
        title: "Failed to delete account",
        message: (error as Error).message,
      });
    }
  };

  const handleDialogOpenChange = (open: boolean) => {
    setIsDialogOpen(open);
    if (!open) {
      setDeleteConfirmText("");
    }
  };

  const isDeleteEnabled = deleteConfirmText === "DELETE";

  if (!currentUser) return null;

  return (
    <div className="space-y-8">
      <Card className="border-destructive/50">
        <CardContent className="p-6">
          <div className="space-y-4">
            <div>
              <h3 className="font-semibold text-destructive">Delete Account</h3>
              <p className="text-sm text-muted-foreground mt-1">
                Permanently delete your account and all associated data. This
                action cannot be undone.
              </p>
            </div>

            <Separator />

            <div className="flex items-center justify-between">
              <div>
                <p className="font-medium">Ready to delete your account?</p>
                <p className="text-sm text-muted-foreground">
                  Once you delete your account, there is no going back. Please
                  be certain.
                </p>
              </div>

              <AlertDialog
                open={isDialogOpen}
                onOpenChange={handleDialogOpenChange}
              >
                <AlertDialogTrigger asChild>
                  <Button variant="destructive" className="ml-4">
                    <Icon type="delete" className="text-sm" />
                    Delete Account
                  </Button>
                </AlertDialogTrigger>
                <AlertDialogContent className="max-w-lg">
                  <AlertDialogHeader>
                    <AlertDialogTitle className="flex items-center gap-2">
                      <Icon type="warning" className="text-destructive" />
                      Delete your account permanently?
                    </AlertDialogTitle>
                    <AlertDialogDescription>
                      This action <strong>cannot be undone</strong>. This will
                      permanently delete your account and remove all your data
                      from our servers.
                    </AlertDialogDescription>
                    <div className="space-y-3">
                      <p className="text-sm font-medium">
                        Please type{" "}
                        <span className="font-mono bg-muted px-1 rounded">
                          DELETE
                        </span>{" "}
                        to confirm:
                      </p>
                      <Input
                        id="delete-confirmation"
                        value={deleteConfirmText}
                        onChange={setDeleteConfirmText}
                        placeholder="Type DELETE to confirm"
                        inputClassName="font-mono"
                      />
                    </div>
                  </AlertDialogHeader>
                  <AlertDialogFooter className="gap-2">
                    <AlertDialogCancel>Cancel</AlertDialogCancel>
                    <AlertDialogAction
                      onClick={handleDeleteAccount}
                      disabled={!isDeleteEnabled}
                      className="bg-destructive hover:bg-destructive/90 disabled:opacity-50"
                    >
                      <Icon type="delete_forever" className="text-sm" />
                      Delete Account Permanently
                    </AlertDialogAction>
                  </AlertDialogFooter>
                </AlertDialogContent>
              </AlertDialog>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
