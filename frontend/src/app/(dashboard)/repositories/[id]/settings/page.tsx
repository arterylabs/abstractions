"use client";

import { isAxiosError } from "axios";
import { FormEvent, useEffect, useMemo, useState } from "react";
import api from "@/api";
import type { TeamMember } from "@/api/team-members";
import useAPI from "@/api/hooks";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { PageHeader } from "@/components/ui/page-header";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import Spinner from "@/components/ui/spinner";
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
} from "@/components/ui/alert-dialog";
import { getErrorMessage } from "@/lib/errors";
import {
  getDayOptions,
  isoWeekdayToName,
  localHourToUtc,
  utcHourToLocal,
} from "@/lib/timezone";
import { useAppDispatch, useAppSelector } from "@/store/hooks";
import {
  deleteRepositoryData,
  updateRepositorySettings,
} from "@/store/modules/repositories";
import { useParams, useRouter } from "next/navigation";

const getLocalHourFromUtc = (utcHour: number) => {
  const now = new Date();
  return new Date(
    Date.UTC(now.getFullYear(), now.getMonth(), now.getDate(), utcHour),
  ).getHours();
};

const HOUR_OPTIONS = Array.from({ length: 24 }, (_, localHour) => ({
  value: String(localHour),
  label: utcHourToLocal(localHourToUtc(localHour)),
}));

export function DeleteRepositoryDataDialog({
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
          <AlertDialogTitle>Reset Repository Data</AlertDialogTitle>
          <AlertDialogDescription>
            This will remove indexed repository data and newsletter history.
            Type
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
            {isDeleting ? "Resetting..." : "Reset Repository Data"}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}

export default function RepositorySettingsPage() {
  const dispatch = useAppDispatch();
  const router = useRouter();
  const params = useParams<{ id: string }>();
  const repoId = params.id;
  const currentWorkspace = useAppSelector(
    (state) => state.auth.currentWorkspace,
  );

  const [repoResponse, loadRepository] = useAPI(api.repositories.get);
  const [membersResponse, loadMembers] = useAPI(api.teamMembers.list);
  const [sendDay, setSendDay] = useState("");
  const [sendHour, setSendHour] = useState("");
  const [newMemberEmail, setNewMemberEmail] = useState("");
  const [memberMessage, setMemberMessage] = useState<string | null>(null);
  const [isSavingSchedule, setIsSavingSchedule] = useState(false);
  const [isAddingMember, setIsAddingMember] = useState(false);
  const [removingMemberId, setRemovingMemberId] = useState<string | null>(null);
  const [isResetDialogOpen, setIsResetDialogOpen] = useState(false);
  const [confirmationValue, setConfirmationValue] = useState("");
  const [isDeleting, setIsDeleting] = useState(false);

  useEffect(() => {
    if (!currentWorkspace?.id || !repoId) {
      return;
    }

    void loadRepository(currentWorkspace.id, repoId);
  }, [currentWorkspace?.id, loadRepository, repoId]);

  useEffect(() => {
    if (!repoId) {
      return;
    }

    void loadMembers(repoId);
  }, [loadMembers, repoId]);

  useEffect(() => {
    if (!repoResponse.result) {
      return;
    }

    setSendDay(String(repoResponse.result.sendDay));
    setSendHour(String(getLocalHourFromUtc(repoResponse.result.sendHour)));
  }, [repoResponse.result]);

  const dayOptions = useMemo(() => getDayOptions(), []);
  const repository = repoResponse.result;
  const members = membersResponse.result ?? [];

  const refreshRepository = async () => {
    if (currentWorkspace?.id) {
      await loadRepository(currentWorkspace.id, repoId);
    }
  };

  const refreshMembers = async () => {
    await loadMembers(repoId);
  };

  const handleSaveSchedule = async () => {
    if (!currentWorkspace?.id || !repository || !sendDay || !sendHour) {
      return;
    }

    setIsSavingSchedule(true);

    try {
      const result = await dispatch(
        updateRepositorySettings({
          workspaceId: currentWorkspace.id,
          repoId: repository.id,
          data: {
            sendDay: Number(sendDay),
            sendHour: localHourToUtc(Number(sendHour)),
          },
        }),
      );

      if ("error" in result) {
        throw new Error(result.payload as string);
      }

      await refreshRepository();
      Notifier.success({
        title: "Schedule updated",
        message: "Repository send schedule saved successfully.",
      });
    } catch (error: unknown) {
      Notifier.error({
        title: "Unable to save schedule",
        message: getErrorMessage(error),
      });
    } finally {
      setIsSavingSchedule(false);
    }
  };

  const handleAddMember = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setMemberMessage(null);
    setIsAddingMember(true);

    try {
      await api.teamMembers.add(repoId, { email: newMemberEmail.trim() });
      setNewMemberEmail("");
      await refreshMembers();
      Notifier.success({
        title: "Team member added",
        message: "They'll receive future newsletters for this repository.",
      });
    } catch (error: unknown) {
      if (isAxiosError(error) && error.response?.status === 402) {
        setMemberMessage("Upgrade to add more members");
      } else {
        setMemberMessage(getErrorMessage(error));
      }
    } finally {
      setIsAddingMember(false);
    }
  };

  const handleRemoveMember = async (member: TeamMember) => {
    setMemberMessage(null);
    setRemovingMemberId(member.id);

    try {
      await api.teamMembers.remove(repoId, member.id);
      await refreshMembers();
      Notifier.success({
        title: "Team member removed",
        message: `${member.email} has been removed from this repository.`,
      });
    } catch (error: unknown) {
      setMemberMessage(getErrorMessage(error));
    } finally {
      setRemovingMemberId(null);
    }
  };

  const handleDeleteData = async () => {
    if (!currentWorkspace?.id || !repository) {
      return;
    }

    setIsDeleting(true);

    try {
      const result = await dispatch(
        deleteRepositoryData({
          workspaceId: currentWorkspace.id,
          repoId: repository.id,
        }),
      );

      if ("error" in result) {
        throw new Error(result.payload as string);
      }

      router.push("/dashboard/repositories");
    } catch (error: unknown) {
      Notifier.error({
        title: "Unable to reset repository",
        message: getErrorMessage(error),
      });
    } finally {
      setIsDeleting(false);
    }
  };

  if (!repository || repoResponse.isFetching) {
    return (
      <div className="flex min-h-[240px] flex-col items-center justify-center gap-4">
        <Spinner />
        <p className="text-sm text-muted-foreground">
          Loading repository settings...
        </p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <PageHeader
        title={repository.fullName}
        showBackButton
        handleBack={() => router.push("/dashboard/repositories")}
      />

      <Card>
        <CardHeader>
          <CardTitle>Send Schedule</CardTitle>
          <CardDescription>
            Choose when your repository newsletter should be sent.
          </CardDescription>
        </CardHeader>
        <CardContent className="grid gap-4 md:grid-cols-2">
          <div className="grid gap-2">
            <label className="text-sm font-medium">Day</label>
            <Select value={sendDay} onValueChange={setSendDay}>
              <SelectTrigger>
                <SelectValue placeholder="Select a day" />
              </SelectTrigger>
              <SelectContent>
                {dayOptions.map((option) => (
                  <SelectItem key={option.value} value={String(option.value)}>
                    {option.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
          <div className="grid gap-2">
            <label className="text-sm font-medium">Time</label>
            <Select value={sendHour} onValueChange={setSendHour}>
              <SelectTrigger>
                <SelectValue placeholder="Select a time" />
              </SelectTrigger>
              <SelectContent>
                {HOUR_OPTIONS.map((option) => (
                  <SelectItem key={option.value} value={option.value}>
                    {option.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
          <div className="md:col-span-2 flex items-center justify-between rounded-lg border p-4">
            <p className="text-sm text-muted-foreground">
              Currently set for {isoWeekdayToName(Number(sendDay))} at{" "}
              {sendHour
                ? HOUR_OPTIONS.find((option) => option.value === sendHour)
                    ?.label
                : "-"}
            </p>
            <Button
              onClick={() => void handleSaveSchedule()}
              loading={isSavingSchedule}
            >
              Save Schedule
            </Button>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Team Members</CardTitle>
          <CardDescription>
            Invite teammates to receive repository newsletters. Free plan: up to
            10 members.
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <form
            className="flex flex-col gap-3 sm:flex-row"
            onSubmit={handleAddMember}
          >
            <Input
              id="team-member-email"
              type="email"
              label="Add member"
              className="flex-1"
              placeholder="teammate@example.com"
              value={newMemberEmail}
              onChange={setNewMemberEmail}
              required
            />
            <Button
              type="submit"
              className="sm:self-end"
              loading={isAddingMember}
            >
              Add Member
            </Button>
          </form>

          {memberMessage && (
            <p className="text-sm text-destructive">{memberMessage}</p>
          )}

          {membersResponse.isFetching && members.length === 0 ? (
            <div className="flex items-center gap-3 text-sm text-muted-foreground">
              <Spinner className="h-5 w-5" />
              Loading team members...
            </div>
          ) : members.length === 0 ? (
            <p className="text-sm text-muted-foreground">
              No team members added yet.
            </p>
          ) : (
            <div className="space-y-3">
              {members.map((member) => (
                <div
                  key={member.id}
                  className="flex flex-col gap-3 rounded-lg border p-4 sm:flex-row sm:items-center sm:justify-between"
                >
                  <div>
                    <p className="font-medium">{member.email}</p>
                    <Badge
                      variant={member.subscribed ? "success" : "secondary"}
                      className="mt-2"
                    >
                      {member.subscribed ? "Subscribed" : "Pending"}
                    </Badge>
                  </div>
                  <Button
                    variant="outline"
                    onClick={() => void handleRemoveMember(member)}
                    loading={removingMemberId === member.id}
                  >
                    Remove
                  </Button>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Data Management</CardTitle>
          <CardDescription>
            Reset indexed repository data and newsletter history. This action
            cannot be undone.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex flex-col gap-4 rounded-lg border border-destructive/30 bg-destructive/5 p-4">
            <p className="text-sm text-muted-foreground">
              Resetting repository data clears indexed code, upcoming topics,
              and stored newsletter data.
            </p>
            <div>
              <Button
                variant="destructive"
                onClick={() => setIsResetDialogOpen(true)}
              >
                Reset Repository Data
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>

      <DeleteRepositoryDataDialog
        open={isResetDialogOpen}
        onOpenChange={setIsResetDialogOpen}
        repositoryFullName={repository.fullName}
        confirmationValue={confirmationValue}
        onConfirmationChange={setConfirmationValue}
        onConfirm={() => void handleDeleteData()}
        isDeleting={isDeleting}
      />
    </div>
  );
}
