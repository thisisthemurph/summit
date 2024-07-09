import { useForm } from 'react-hook-form';
import { FormField } from "../../shared/components/Forms.tsx";
import { z, ZodType } from "zod";
import { zodResolver } from "@hookform/resolvers/zod"

type FormValues = {
  email: string;
  password: string;
}

const LoginSchema: ZodType<FormValues> = z
  .object({
    email: z.string().email("Email is not valid"),
    password: z.string()
      .min(6, "Password must be at least 6 characters long")
      .max(24, "Password cannot be more than 24 characters long"),
  });

function LoginPage() {
  const {
    register,
    handleSubmit,
    formState: {errors}
  } = useForm<FormValues>({ resolver: zodResolver(LoginSchema)});

  const onSubmit = handleSubmit(async (data) => {
    const baseUrl = import.meta.env.VITE_API_BASE_URL;
    const result = await fetch(`${baseUrl}/login`, {method: "POST", body: JSON.stringify(data)});
    console.log(result);
  })

  return (
    <>
      <section className="prose mb-6">
        <h1>Login</h1>
      </section>
      <form onSubmit={onSubmit} className="space-y-4">
        <FormField type="email" label="Email address" placeholder="Your email address" register={register("email")} error={errors.email} />
        <FormField type="password" label="Password" placeholder="Your password" register={register("password")} error={errors.password} />
        <button type="submit" className="btn btn-primary">Login</button>
      </form>
    </>
  )
}

export default LoginPage;