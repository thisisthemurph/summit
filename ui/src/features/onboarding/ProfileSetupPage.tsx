import {FormField} from "../../shared/components/Forms.tsx";
import {z, ZodType} from "zod";
import {useForm} from "react-hook-form";
import {zodResolver} from "@hookform/resolvers/zod";
import {useEffect, useState} from "react";

type FormValues = {
  firstName: string;
  lastName: string;
}

const ProfileSchema: ZodType<FormValues> = z
  .object({
    firstName: z.string()
      .min(2, "Your first name must be at least 2 characters")
      .max(50, "Your first name cannot be more than 50 characters"),
    lastName: z.string()
      .min(2, "Your last name must be at least 2 characters")
      .max(50, "Your last name cannot be more than 50 characters"),
  });

function ProfileSetupPage() {
  const [isLoading, setIsLoading] = useState(true);
  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<FormValues>({ resolver: zodResolver(ProfileSchema) });

  useEffect(() => {
    const getProfileData = async () => {
      try {
        const baseUrl = import.meta.env.VITE_API_BASE_URL;
        const response = await fetch(`${baseUrl}/onboarding/profile`, {
          method: "GET",
          credentials: "include",
          headers: {
            "Content-Type": "application/json",
          }
        });
        const data = await response.json();
        reset(data);
      } catch {
        console.error("Error fetching profile data");
      } finally {
        setIsLoading(false);
      }
    }

    getProfileData();
  }, [reset]);

  const onSubmit = handleSubmit(async (data) => {
    console.log(data);
    const baseUrl = import.meta.env.VITE_API_BASE_URL;
    const result = await fetch(`${baseUrl}/onboarding/profile`, {
      method: "POST",
      body: JSON.stringify(data),
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      }
    });

    if (result.status !== 200) {
      alert("There has been an issue updating your personal data!");
    }
  })

  return (
    <>
      <section className="prose mb-6">
        <h1>Set up your profile</h1>
      </section>
      <form onSubmit={onSubmit} className="space-y-4">
        <FormField label="First name" register={register("firstName")} error={errors.firstName}/>
        <FormField label="Last name" register={register("lastName")} error={errors.lastName}/>
        <button type="submit" className="btn btn-primary" disabled={isLoading}>Continue</button>
      </form>
    </>
  )
}

export default ProfileSetupPage;