import { useState } from "react";
import { Input } from "../ui/input";
import { DataTable, DataTableColumn } from "../ui/data-table";
import { Button } from "../ui/button";

interface User {
  id: string;
  name: string;
  email: string;
}

interface UsersTableProps {
  users: User[];
  onSearch: (query: string) => void;
  onEdit: (user: User) => void;
  onDelete: (userId: string) => void;
  onViewDetails: (user: User) => void;
}

export default function UsersTable({
  users,
  onSearch,
  onEdit,
  onDelete,
  onViewDetails,
}: UsersTableProps) {
  const [searchQuery, setSearchQuery] = useState("");

  const handleSearch = (value: string) => {
    setSearchQuery(value);
    onSearch(value);
  };

  const columns: DataTableColumn<User>[] = [
    {
      key: "name",
      label: "Name",
      order: 1,
    },
    {
      key: "email",
      label: "Email",
      order: 2,
    },
    {
      key: "actions",
      label: "Actions",
      order: 3,
      render: (_, user) => (
        <div className="flex gap-2">
          <Button onClick={() => onViewDetails(user)}>View</Button>
          <Button onClick={() => onEdit(user)}>Edit</Button>
          <Button color="red" onClick={() => onDelete(user.id)}>
            Delete
          </Button>
        </div>
      ),
    },
  ];

  return (
    <div>
      <div className="mb-4">
        <Input
          label="Search"
          id="search"
          placeholder="Search users..."
          value={searchQuery}
          onChange={(v) => handleSearch(v)}
        />
      </div>

      <DataTable data={users} columns={columns} />
    </div>
  );
}
