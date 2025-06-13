// @ts-check
const path = require('path');
const grafanaConfig = require('@grafana/eslint-config/flat');
const grafanaI18nPlugin = require('@grafana/i18n/eslint-plugin')
const importPlugin = require('eslint-plugin-import');
const jsxA11yPlugin = require('eslint-plugin-jsx-a11y');
const jestPlugin = require('eslint-plugin-jest');
const { includeIgnoreFile } = require('@eslint/compat');

/**
 * @type {Array<import('eslint').Linter.Config>}
 */
module.exports = [
  includeIgnoreFile(path.resolve(__dirname, '.gitignore')),
  grafanaConfig,
  {
    plugins: {
      'import': importPlugin,
    },
    files: ['src/**/*.{ts,tsx}'],
    rules: {
      '@typescript-eslint/no-deprecated': 'warn',
      'import/order': [
        'error',
        {
          'groups': [['builtin', 'external'], 'internal', 'parent', 'sibling', 'index'],
          'newlines-between': 'always',
          'alphabetize': { 'order': 'asc' }
        }
      ],
    },
    settings: {
      'import/external-module-folders': ['node_modules', '.yarn'],
    },
    languageOptions: {
      parserOptions: {
        project: './tsconfig.json'
      }
    }
  },
  {
    name: 'test-overrides',
    files: ['src/**/*.{test,spec}.{ts,tsx}'],
    plugins: {
      jest: jestPlugin
    },
    rules: {
      'jest/no-focused-tests': 'error',
    }
  },
  {
    // Sections of codebase that have all translation markup issues fixed
    name: 'grafana/i18n-overrides',
    plugins: {
      '@grafana/i18n': grafanaI18nPlugin,
    },
    files: [
      'src/**/*.{ts,tsx,js,jsx}',
    ],
    ignores: [
      'src/**/*.{test,spec,story}.{ts,tsx}',
    ],
    rules: {
      '@grafana/i18n/no-untranslated-strings': ['error', { calleesToIgnore: ['^css$', 'use[A-Z].*'] }],
      '@grafana/i18n/no-translation-top-level': 'error',
    },
  },
  {
    name: 'a11y-overrides',
    files: ['**/*.tsx'],
    ignores: ['**/*.{spec,test}.tsx'],
    plugins: {
      'jsx-a11y': jsxA11yPlugin
    },
    rules: {
      ...jsxA11yPlugin.configs.recommended.rules,
      // rules marked "off" are those left in the recommended preset we need to fix
      // we should remove the corresponding line and fix them one by one
      // any marked "error" contain specific overrides we'll need to keep
      'jsx-a11y/no-autofocus': [
        'error',
        {
          ignoreNonDOM: true
        }
      ],
      'jsx-a11y/label-has-associated-control': [
        'error',
        {
          controlComponents: ['NumberInput'],
          depth: 2
        }
      ]
    }
  }
];
