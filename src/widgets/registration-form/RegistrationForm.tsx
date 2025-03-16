'use client';

import { AppConfig } from 'next-intl';
import { FC } from 'react';

import { Button, TextField } from 'shared/ui';
import { IssueCode, clsx, formResolver, useForm, validate } from 'shared/lib';

interface RegistrationFormProps {
  className?: string;
  messages?: {
    auth?: AppConfig['Messages']['auth'];
    validation?: AppConfig['Messages']['validation'];
  };
}

export const RegistrationForm: FC<RegistrationFormProps> = ({
  className,
  messages,
}) => {
  const schema = validate
    .object({
      email: validate
        .string()
        .nonempty(messages?.validation?.['validation-required'])
        .email(messages?.validation?.['validation-email']),
      username: validate
        .string()
        .nonempty(messages?.validation?.['validation-required'])
        .min(4, messages?.validation?.['validation-username-length'])
        .max(30, messages?.validation?.['validation-username-length'])
        .regex(
          /^[a-zA-z]+[0-9]*[@-_]?$/,
          messages?.validation?.['validation-username'],
        ),
      password: validate
        .string()
        .nonempty(messages?.validation?.['validation-required'])
        .min(6, messages?.validation?.['validation-password-length'])
        .max(75, messages?.validation?.['validation-password-length'])
        .regex(
          /^[a-zA-Z]+[0-9]+[@-_]?$/,
          messages?.validation?.['validation-password-symbols'],
        ),
      confirmPassword: validate
        .string()
        .nonempty(messages?.validation?.['validation-required']),
    })
    .superRefine((fields, ctx) => {
      if (
        fields.confirmPassword &&
        fields.password !== fields.confirmPassword
      ) {
        ctx.addIssue({
          path: ['confirmPassword'],
          message: messages?.validation?.['validation-match-passwords'],
          code: IssueCode.custom,
        });
      }
    });
  const {
    register,
    handleSubmit,
    formState: { isSubmitting, errors },
  } = useForm({
    resolver: formResolver(schema),
  });

  const onSubmit = handleSubmit(async () => {});

  return (
    <div
      className={clsx(
        className,
        `
        flex flex-col gap-xl bg-surface p-l rounded-m shadow-surface tablet:max-w-[630px]
        `,
      )}
    >
      <h4 className="text-xl font-bold whitespace-pre-wrap tablet:text-xxl">
        {messages?.auth?.['registration-tagline']}
      </h4>
      <form
        className="flex flex-col gap-2xl"
        onSubmit={onSubmit}
      >
        <div className="flex flex-col gap-m">
          <TextField
            {...register('email')}
            autoComplete="email"
            placeholder={messages?.auth?.['email-input']}
            errorMessage={errors.email?.message}
          />
          <TextField
            {...register('username')}
            autoComplete="username"
            placeholder={messages?.auth?.['username-input']}
            errorMessage={errors.username?.message}
          />
          <TextField
            {...register('password')}
            type="password"
            autoComplete="new-password"
            placeholder={messages?.auth?.['password-input']}
            errorMessage={errors.password?.message}
          />
          <TextField
            {...register('confirmPassword')}
            type="password"
            placeholder={messages?.auth?.['config-password-input']}
            errorMessage={errors.confirmPassword?.message}
          />
        </div>
        <Button
          type="submit"
          disabled={isSubmitting}
          loading={isSubmitting}
          variant="primaryFilled"
        >
          {messages?.auth?.['registration-button']}
        </Button>
      </form>
    </div>
  );
};
