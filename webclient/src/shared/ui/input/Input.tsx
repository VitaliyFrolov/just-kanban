import { FC, InputHTMLAttributes, ReactNode, Ref } from 'react';

import { tv } from 'shared/lib';

export const inputVariants = tv({
  slots: {
    root: `
      flex border-1 border-contour rounded-md overflow-hidden
      duration-[var(--duration-interactive)] focus-within:border-accented
    `,
    nativeInput:
      'px-[15px] py-[12px] flex-auto bg-transparent focus-visible:outline-none placeholder:text-placeholder min-w-[100px]',
  },
  variants: {
    error: {
      true: {
        root: 'border-danger-500 text-danger-500',
        // nativeInput: 'placeholder:text-danger-500',
      },
    },
  },
});

export interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  icon?: ReactNode;
  error?: boolean;
  beforeInput?: ReactNode;
  afterInput?: ReactNode;
  ref?: Ref<HTMLInputElement>;
}

export const Input: FC<InputProps> = ({
  className,
  icon,
  error,
  afterInput,
  beforeInput,
  autoComplete = 'off',
  ref,
  ...props
}) => {
  const { root, nativeInput } = inputVariants({ error });

  return (
    <div className={root({ className })}>
      {beforeInput}
      <div className="w-[43px] flex shrink-0 justify-end self-center h-full relative empty:hidden">
        {icon && (
          <span className="w-[25px] height-[25px] flex items-center justify-center">
            {icon}
          </span>
        )}
      </div>
      <input
        {...props}
        ref={ref}
        autoComplete={autoComplete}
        className={nativeInput()}
      />
      {afterInput}
    </div>
  );
};
