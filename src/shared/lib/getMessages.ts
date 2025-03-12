'use server';

import { AppConfig, createTranslator } from 'next-intl';
import { ReactNode } from 'react';
import { getLocale } from 'next-intl/server';

import { promises as fs } from 'fs';
import path from 'path';

type CustomTranslationFunction<NS extends keyof AppConfig['Messages']> = {
  (key: keyof AppConfig['Messages'][NS]): string;
  rich: (
    key: keyof AppConfig['Messages'][NS],
    options?: Record<string, string | number | ((chunk: string) => ReactNode)>,
  ) => string | ReactNode;
};

export async function getMessages<NS extends keyof AppConfig['Messages']>(
  namespace: NS,
) {
  const locale = await getLocale();
  const filePath = path.join(
    process.cwd(),
    'public',
    'messages',
    locale,
    `${namespace}.json`,
  );
  const fileContent = await fs.readFile(filePath, 'utf-8');
  const rawMessages = JSON.parse(fileContent);
  const messages = { [namespace]: rawMessages };

  const t = createTranslator({
    locale,
    namespace,
    messages,
  });

  return t as unknown as CustomTranslationFunction<NS>;
}
