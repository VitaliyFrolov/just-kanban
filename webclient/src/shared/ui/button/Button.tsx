import { ButtonHTMLAttributes, FC, ReactNode, Ref } from 'react';

import { VariantProps, tv } from 'shared/lib';
import { Loader } from 'shared/ui';

export const buttonVariants = tv({
  base: `
    border-1 inline-flex rounded-md disabled:opacity-60 relative
    justify-center duration-[var(--duration-interactive)]
  `,
  variants: {
    variant: {
      clear: 'bg-transparent p-0 text-inherit',
      outlined: 'border-1 border-primary bg-primary',
      primaryFilled:
        'bg-primary-900 text-white font-bold hover:bg-primary-950 focus-visible:bg-primary-950',
    },
    size: {
      md: 'py-m px-xl',
      xl: 'py-l px-xxl',
    },
    icons: {
      true: 'px-[22px] flex items-center justify-center gap-[10px]',
      false: '',
    },
  },
  defaultVariants: {
    variant: 'outlined',
    size: 'md',
  },
});

interface ButtonProps
  extends ButtonHTMLAttributes<HTMLButtonElement>,
    VariantProps<typeof buttonVariants> {
  leftIcon?: ReactNode;
  rightIcon?: ReactNode;
  ref?: Ref<HTMLButtonElement>;
  loading?: boolean;
}

export const Button: FC<ButtonProps> = ({
  className,
  size,
  variant,
  leftIcon,
  rightIcon,
  ref,
  loading,
  type = 'button',
  ...props
}) => {
  return (
    <button
      {...props}
      ref={ref}
      className={buttonVariants({
        className,
        variant,
        size,
        icons: !!leftIcon || !!rightIcon,
      })}
      type={type}
    >
      {loading && (
        <div className="bg-inherit absolute inset-0 flex items-center justify-center">
          <Loader />
        </div>
      )}
      {leftIcon && (
        <span className="flex w-[25px] h-[25px] items-center justify-center">
          {leftIcon}
        </span>
      )}
      {props.children}
      {rightIcon && (
        <span className="flex w-[25px] h-[25px] items-center justify-center">
          {rightIcon}
        </span>
      )}
    </button>
  );
};
