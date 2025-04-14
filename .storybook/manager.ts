import { addons } from '@storybook/manager-api';
import { create } from '@storybook/theming/create';

const theme = create({
  base: 'light',
  brandTitle: 'My UI Component Library',
  brandUrl: 'https://example.com',
  brandTarget: '_self',
  
  // UI colors
  colorPrimary: '#1ea7fd',
  colorSecondary: '#ff4785',

  // UI
  appBg: '#f6f9fc',
  appContentBg: '#ffffff',
  appBorderColor: '#e6e8eb',
  appBorderRadius: 4,

  // Typography
  fontBase: '"Open Sans", sans-serif',
  fontCode: 'monospace',

  // Text colors
  textColor: '#333333',
  textInverseColor: '#ffffff',
});

addons.setConfig({
  theme,
  showToolbar: true,
}); 