import { FlatCompat } from '@eslint/eslintrc';
import { fileURLToPath } from 'node:url';
import js from '@eslint/js';
import path from 'node:path';
import tsParser from '@typescript-eslint/parser';
import typescriptEslint from '@typescript-eslint/eslint-plugin';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const compat = new FlatCompat({
  baseDirectory: __dirname,
  recommendedConfig: js.configs.recommended,
  allConfig: js.configs.all,
});

export default [
  {
    ignores: ['**/node_modules', '**/.next'],
  },
  ...compat.extends(
    'next/core-web-vitals',
    'plugin:@typescript-eslint/eslint-recommended',
    'plugin:@typescript-eslint/recommended',
    'plugin:prettier/recommended',
  ),
  {
    plugins: {
      '@typescript-eslint': typescriptEslint,
    },

    languageOptions: {
      parser: tsParser,
    },

    rules: {
      'prettier/prettier': [
        'error',
        {
          singleQuote: true,
          jsxSingleQuote: false,
          printWidth: 80,
          tabWidth: 2,
          useTabs: false,
          semi: true,
          quoteProps: 'as-needed',
          trailingComma: 'all',
          bracketSpacing: true,
          bracketSameLine: false,
          arrowParens: 'always',
          endOfLine: 'lf',
          singleAttributePerLine: true,
        },
      ],

      'no-console': 'error',
      '@typescript-eslint/no-unused-vars': 'error',
      'react/display-name': 'error',
      'import/no-anonymous-default-export': 'off',
      'sort-imports': [
        'error',
        {
          allowSeparatedGroups: true,
        },
      ],

      'react/jsx-curly-brace-presence': ['error', 'never'],

      'react/function-component-definition': [
        'error',
        {
          namedComponents: 'arrow-function',
          unnamedComponents: 'arrow-function',
        },
      ],
    },
  },
];
