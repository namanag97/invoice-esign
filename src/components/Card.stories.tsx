import React from 'react';
import type { Meta, StoryObj } from '@storybook/react';
import { Card } from './Card';

const meta = {
  title: 'Components/Card',
  component: Card,
  parameters: {
    layout: 'centered',
  },
  tags: ['autodocs'],
  argTypes: {
    borderColor: { control: 'color' },
    width: { control: 'text' },
  },
} satisfies Meta<typeof Card>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    title: 'Card Title',
    children: (
      <p>This is a simple card component with some content. It can be used to display information in a contained area.</p>
    ),
  },
};

export const Wide: Story = {
  args: {
    title: 'Wide Card',
    width: '500px',
    children: (
      <div>
        <p>This is a wider card with more content.</p>
        <p>You can adjust the width to fit your needs.</p>
      </div>
    ),
  },
};

export const Colored: Story = {
  args: {
    title: 'Colored Border Card',
    borderColor: '#1ea7fd',
    children: (
      <p>This card has a custom border color that can be set through props.</p>
    ),
  },
}; 