"use client";

import { Input } from "@/components/ui/input";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";

export function DeleteRepositoryDialog({
  open,
  onOpenChange,
  repositoryFullName,
  confirmationValue,
  onConfirmationChange,
  onConfirm,
  isDeleting,
}: {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  repositoryFullName: string;
  confirmationValue: string;
  onConfirmationChange: (value: string) => void;
  onConfirm: () => void;
  isDeleting: boolean;
}) {
  const isMatch = confirmationValue === repositoryFullName;

  return (
    <AlertDialog open={open} onOpenChange={onOpenChange}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Delete Repository</AlertDialogTitle>
          <AlertDialogDescription>
            This will permanently delete the repository and all its data. This
            action cannot be undone. Type
            <span className="font-medium text-foreground">
              {" "}
              {repositoryFullName}
            </span>{" "}
            to confirm.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <Input
          id="confirm-repository-name"
          label="Repository name"
          value={confirmationValue}
          onChange={onConfirmationChange}
          placeholder={repositoryFullName}
        />
        <AlertDialogFooter>
          <AlertDialogCancel>Cancel</AlertDialogCancel>
          <AlertDialogAction
            onClick={onConfirm}
            disabled={!isMatch || isDeleting}
            className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
          >
            {isDeleting ? "Deleting..." : "Delete Repository"}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}
