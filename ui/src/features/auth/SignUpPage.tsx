import { useForm } from 'react-hook-form';
import { FormField } from "../../shared/components/Forms.tsx";
import { z, ZodType } from "zod";
import { zodResolver } from "@hookform/resolvers/zod"
import axiosInstance from "../../shared/requests/axiosInstance";
import {AxiosError} from "axios";
import Container from "../../shared/components/Container.tsx";

type FormValues = {
  email: string;
  password: string;
  confirmPassword: string;
}

type UserResponse = {
  id: string;
  firstName: string;
  lastName: string;
  email: string;
}

const SignUpSchema: ZodType<FormValues> = z
  .object({
    email: z.string().email("Email is not valid"),
    password: z.string()
      // .min(6, "Password must be at least 6 characters long")
      .max(24, "Password cannot be more than 24 characters long"),
    confirmPassword: z.string(),
  }).refine((data) => data.password === data.confirmPassword, {
    message: "Passwords do not match",
    path: ["confirmPassword"],
  });

function SignUpPage() {
  const {
    register,
    handleSubmit,
    formState: {errors}
  } = useForm<FormValues>({ resolver: zodResolver(SignUpSchema)});

  const onSubmit = handleSubmit(async (data) => {
    try {
      const response = await axiosInstance.post<UserResponse>("/signup", data);
      console.log(response.data);
    } catch (e) {
      let message = "There has been an error signing you up";
      if (e instanceof AxiosError) {
        message = e.response?.data?.message ? e.response.data.message : message;
      }
      alert(message);
    }
  });

  return (
    <Container>
      <form onSubmit={onSubmit} className="space-y-4">
        <h1 className="text-2xl">Sign up</h1>
        <FormField type="email"
                   label="Email address"
                   value={import.meta.env.VITE_DEV_EMAIL}
                   placeholder="Your email address"
                   register={register("email")}
                   error={errors.email}/>
        <FormField type="password"
                   label="Password"
                   placeholder="Your password"
                   register={register("password")}
                   error={errors.password}/>
        <FormField type="password"
                   label="Confirm password"
                   placeholder="Confirm your password"
                   register={register("confirmPassword")}
                   error={errors.confirmPassword}/>
        <button type="submit" className="btn btn-primary">Signup</button>
      </form>
    </Container>
  )
}

export default SignUpPage;