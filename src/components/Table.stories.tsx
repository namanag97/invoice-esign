import React from 'react';
import type { Meta, StoryObj } from '@storybook/react';
import { Table } from './Table';

// Mock data for the table
interface User {
  id: number;
  name: string;
  email: string;
  role: string;
  status: 'active' | 'inactive' | 'pending';
  lastLogin: string;
}

const users: User[] = [
  {
    id: 1,
    name: 'Alex Johnson',
    email: 'alex@example.com',
    role: 'Admin',
    status: 'active',
    lastLogin: '2023-11-10T14:30:00'
  },
  {
    id: 2,
    name: 'Sarah Williams',
    email: 'sarah@example.com',
    role: 'User',
    status: 'active',
    lastLogin: '2023-11-09T09:15:00'
  },
  {
    id: 3,
    name: 'Michael Brown',
    email: 'michael@example.com',
    role: 'Editor',
    status: 'inactive',
    lastLogin: '2023-10-25T11:45:00'
  },
  {
    id: 4,
    name: 'Emily Davis',
    email: 'emily@example.com',
    role: 'User',
    status: 'pending',
    lastLogin: '2023-11-08T16:20:00'
  },
  {
    id: 5,
    name: 'James Wilson',
    email: 'james@example.com',
    role: 'Viewer',
    status: 'active',
    lastLogin: '2023-11-10T10:05:00'
  }
];

// Format the date for display
const formatDate = (dateString: string) => {
  const date = new Date(dateString);
  return new Intl.DateTimeFormat('en-US', {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  }).format(date);
};

// Status badge component
const StatusBadge = ({ status }: { status: User['status'] }) => {
  const getStatusStyles = () => {
    switch (status) {
      case 'active':
        return { backgroundColor: '#dcfce7', color: '#15803d' };
      case 'inactive':
        return { backgroundColor: '#fef2f2', color: '#b91c1c' };
      case 'pending':
        return { backgroundColor: '#fff7ed', color: '#c2410c' };
      default:
        return { backgroundColor: '#f3f4f6', color: '#4b5563' };
    }
  };

  return (
    <span
      style={{
        ...getStatusStyles(),
        padding: '2px 8px',
        borderRadius: '12px',
        fontSize: '12px',
        fontWeight: 500,
        display: 'inline-block',
        textTransform: 'capitalize'
      }}
    >
      {status}
    </span>
  );
};

const meta: Meta<typeof Table> = {
  title: 'Components/Table',
  component: Table,
  parameters: {
    layout: 'padded',
  },
  tags: ['autodocs'],
};

export default meta;
type Story = StoryObj<typeof Table>;

export const Default: Story = {
  args: {
    columns: [
      { header: 'Name', accessor: 'name' },
      { header: 'Email', accessor: 'email' },
      { header: 'Role', accessor: 'role' },
      { 
        header: 'Status', 
        accessor: (user) => <StatusBadge status={user.status} /> 
      },
      { 
        header: 'Last Login', 
        accessor: (user) => formatDate(user.lastLogin) 
      },
    ],
    data: users,
    hoverEffect: true,
  },
};

export const WithRowClick: Story = {
  args: {
    ...Default.args,
    onRowClick: (item) => alert(`Row clicked: ${item.name}`),
  },
};

export const FixedHeader: Story = {
  args: {
    ...Default.args,
    fixedHeader: true,
    maxHeight: '300px',
  },
};

export const EmptyTable: Story = {
  args: {
    ...Default.args,
    data: [],
  },
}; 