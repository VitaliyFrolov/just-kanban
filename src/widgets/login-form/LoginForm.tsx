'use client';

import { AppConfig } from 'next-intl';
import { FC } from 'react';

import { Button, TextField } from 'shared/ui';
import { clsx, formResolver, useForm, validate } from 'shared/lib';

interface LoginFormProps {
  className?: string;
  messages?: {
    auth?: AppConfig['Messages']['auth'];
    validation?: AppConfig['Messages']['validation'];
  };
}

export const LoginForm: FC<LoginFormProps> = ({ className, messages }) => {
  const schema = validate.object({
    identifier: validate
      .string()
      .nonempty(messages?.validation?.['validation-required']),
    password: validate
      .string()
      .nonempty(messages?.validation?.['validation-required']),
    remember: validate.boolean(),
  });
  const {
    register,
    handleSubmit,
    formState: { isSubmitting, errors },
  } = useForm({
    defaultValues: {
      remember: false,
    },
    resolver: formResolver(schema),
  });

  const onSubmit = handleSubmit(async (fields) => {
    console.log(fields);
  });

  return (
    <div
      className={clsx(
        className,
        `
        flex flex-col gap-xl bg-surface p-l rounded-m shadow-surface
        tablet:max-w-[630px]
        `,
      )}
    >
      <h4 className="text-xl font-bold whitespace-pre-wrap tablet:text-xxl">
        {messages?.auth?.['sign-in-tagline']}
      </h4>
      <form
        className="flex flex-col gap-l"
        onSubmit={onSubmit}
      >
        <TextField
          {...register('identifier')}
          autoComplete="username webauthn"
          placeholder={messages?.auth?.['identifier-input']}
          errorMessage={errors.identifier?.message}
        />
        <TextField
          {...register('password')}
          type="password"
          autoComplete="current-password"
          placeholder={messages?.auth?.['password-input']}
          errorMessage={errors.password?.message}
        />
        <Button
          type="submit"
          disabled={isSubmitting}
          loading={isSubmitting}
          variant="primaryFilled"
        >
          {messages?.auth?.['login-button']}
        </Button>
      </form>
    </div>
  );
};
