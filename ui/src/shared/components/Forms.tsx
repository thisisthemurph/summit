type InputType = "text" | "email" | "number" | "password";

import {FieldError, UseFormRegisterReturn} from "react-hook-form";

type InputProps = {
  type?: InputType;
  label: string;
  placeholder?: string;
  register: UseFormRegisterReturn;
  error: FieldError | undefined;
}

const FormField = ({type, label, placeholder, register, error}: InputProps) => {
  return (
    <label className="form-control w-full">
      <div className="label">
        <span className="label-text">{ label }</span>
      </div>
      <input
        {...register}
        type={type ?? "text"}
        placeholder={placeholder}
        className="input input-bordered w-full"
      />
      {error && (
        <label className="text-sm text-error">{error.message}</label>
      )}
    </label>
  )
}

export {FormField}
