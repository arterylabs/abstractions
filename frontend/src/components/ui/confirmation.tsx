"use client";

import * as React from "react";
import { createPortal } from "react-dom";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "./alert-dialog";
import { cn } from "@/lib/utils";

interface ConfirmationOptions {
  title: string;
  content: string;
  variant?: "primary" | "destructive";
  okText?: string;
  cancelText?: string;
  okCallback: () => void | Promise<void>;
  cancelCallback?: () => void | Promise<void>;
}

interface ConfirmationDialogProps extends ConfirmationOptions {
  open: boolean;
  onClose: () => void;
}

function ConfirmationDialog({
  title,
  content,
  variant = "primary",
  okText = "Confirm",
  cancelText = "Cancel",
  okCallback,
  cancelCallback,
  open,
  onClose,
}: ConfirmationDialogProps) {
  const [isLoading, setIsLoading] = React.useState(false);

  const handleOk = async () => {
    setIsLoading(true);
    try {
      await okCallback();
      onClose();
    } catch (error) {
      // Error handling can be added here if needed
      console.error("Confirmation callback error:", error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleCancel = async () => {
    if (cancelCallback) {
      try {
        await cancelCallback();
      } catch (error) {
        console.error("Cancel callback error:", error);
      }
    }
    onClose();
  };

  return (
    <AlertDialog open={open} onOpenChange={onClose}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>{title}</AlertDialogTitle>
          <AlertDialogDescription>{content}</AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel onClick={handleCancel} disabled={isLoading}>
            {cancelText}
          </AlertDialogCancel>
          <AlertDialogAction
            onClick={handleOk}
            disabled={isLoading}
            className={cn(
              variant === "destructive" &&
                "bg-destructive hover:bg-destructive/90 text-destructive-foreground",
            )}
          >
            {isLoading ? "Processing..." : okText}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}

// Portal wrapper component to manage the confirmation dialog
function ConfirmationPortal({
  title,
  content,
  variant,
  okText,
  cancelText,
  okCallback,
  cancelCallback,
  onClose,
}: ConfirmationOptions & { onClose: () => void }) {
  const [isOpen, setIsOpen] = React.useState(true);
  const [portalRoot, setPortalRoot] = React.useState<HTMLElement | null>(null);

  React.useEffect(() => {
    // Get the portal root element from the layout
    const root = document.getElementById("portal-root");
    setPortalRoot(root);
  }, []);

  const handleClose = () => {
    setIsOpen(false);
    onClose();
  };

  if (!isOpen || !portalRoot) {
    return null;
  }

  return createPortal(
    <ConfirmationDialog
      title={title}
      content={content}
      variant={variant}
      okText={okText}
      cancelText={cancelText}
      okCallback={okCallback}
      cancelCallback={cancelCallback}
      open={isOpen}
      onClose={handleClose}
    />,
    portalRoot,
  );
}

// Global state management for confirmation dialogs
let confirmationQueue: Array<{
  id: string;
  options: ConfirmationOptions;
  cleanup: () => void;
}> = [];

let setGlobalConfirmations: React.Dispatch<
  React.SetStateAction<typeof confirmationQueue>
> | null = null;

// Global confirmation manager component
function GlobalConfirmationManager() {
  const [confirmations, setConfirmations] = React.useState(confirmationQueue);

  React.useEffect(() => {
    setGlobalConfirmations = setConfirmations;
    return () => {
      setGlobalConfirmations = null;
    };
  }, []);

  return (
    <>
      {confirmations.map((confirmation) => (
        <ConfirmationPortal
          key={confirmation.id}
          title={confirmation.options.title}
          content={confirmation.options.content}
          variant={confirmation.options.variant}
          okText={confirmation.options.okText}
          cancelText={confirmation.options.cancelText}
          okCallback={confirmation.options.okCallback}
          cancelCallback={confirmation.options.cancelCallback}
          onClose={confirmation.cleanup}
        />
      ))}
    </>
  );
}

/**
 * Shows a confirmation dialog with customizable title, content, and actions.
 *
 * @param options - Configuration object for the confirmation dialog
 * @param options.title - The title of the confirmation dialog
 * @param options.content - The content/description text of the dialog
 * @param options.variant - The visual style variant: "primary" (default) or "destructive"
 * @param options.okText - Text for the confirmation button (default: "Confirm")
 * @param options.cancelText - Text for the cancel button (default: "Cancel")
 * @param options.okCallback - Function to execute when user confirms (required)
 * @param options.cancelCallback - Optional function to execute when user cancels
 *
 * @example
 * ```tsx
 * import { askConfirmation } from "@/components/ui/confirmation";
 *
 * const handleDelete = () => {
 *   askConfirmation({
 *     title: "Delete Item",
 *     content: "Are you sure you want to delete this item? This action cannot be undone.",
 *     variant: "destructive",
 *     okText: "Delete",
 *     cancelText: "Cancel",
 *     okCallback: async () => {
 *       await deleteItem(id);
 *       // Item deleted successfully
 *     }
 *   });
 * };
 * ```
 */
export function askConfirmation(options: ConfirmationOptions): void;
export function askConfirmation(
  title: string,
  content: string,
  variant?: "primary" | "destructive",
  okText?: string,
  cancelText?: string,
  okCallback?: () => void | Promise<void>,
  cancelCallback?: () => void | Promise<void>,
): void;
export function askConfirmation(
  titleOrOptions: string | ConfirmationOptions,
  content?: string,
  variant: "primary" | "destructive" = "primary",
  okText: string = "Confirm",
  cancelText: string = "Cancel",
  okCallback?: () => void | Promise<void>,
  cancelCallback?: () => void | Promise<void>,
): void {
  let options: ConfirmationOptions;

  // Handle both parameter styles for backward compatibility
  if (typeof titleOrOptions === "object") {
    options = {
      variant: "primary",
      okText: "Confirm",
      cancelText: "Cancel",
      ...titleOrOptions,
    };
  } else {
    if (!okCallback) {
      throw new Error(
        "okCallback is required when using positional parameters",
      );
    }
    options = {
      title: titleOrOptions,
      content: content!,
      variant,
      okText,
      cancelText,
      okCallback,
      cancelCallback,
    };
  }

  const id = Math.random().toString(36).substr(2, 9);

  const cleanup = () => {
    confirmationQueue = confirmationQueue.filter((c) => c.id !== id);
    if (setGlobalConfirmations) {
      setGlobalConfirmations([...confirmationQueue]);
    }
  };

  const confirmation = {
    id,
    options,
    cleanup,
  };

  confirmationQueue.push(confirmation);

  if (setGlobalConfirmations) {
    setGlobalConfirmations([...confirmationQueue]);
  }
}

// Export the dialog component as well for cases where state management is needed
export { ConfirmationDialog, GlobalConfirmationManager };
