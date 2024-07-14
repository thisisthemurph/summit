import {FormField} from "../../shared/components/Forms";
import {z, ZodType} from "zod";
import {useForm} from "react-hook-form";
import {zodResolver} from "@hookform/resolvers/zod";
import {useEffect, useState} from "react";
import axiosInstance from "../../shared/requests/axiosInstance";
import PageHeader from "../../shared/components/PageHeader.tsx";
import Container from "../../shared/components/Container.tsx";
import {useNavigate} from "react-router-dom";


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
  const navigate = useNavigate();
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
        const response = await axiosInstance.get<FormValues>("/onboarding/profile");
        reset(response.data);
      } finally {
        setIsLoading(false);
      }
    }

    setIsLoading(true);
    getProfileData();
  }, [reset]);

  const onSubmit = handleSubmit(async (data) => {
    try {
      await axiosInstance.post("/onboarding/profile", data);
      navigate("/dashboard");
    } catch {
      alert("There has been an issue updating your personal data!");
    }
  })

  return (
    <>
      <PageHeader title="Complete your profile" subtitle="Complete the onboarding process before continuing..." />
      <Container>
        <form onSubmit={onSubmit} className="space-y-4">
          <FormField label="First name" register={register("firstName")} error={errors.firstName}/>
          <FormField label="Last name" register={register("lastName")} error={errors.lastName}/>
          <button type="submit" className="btn btn-primary" disabled={isLoading}>Continue</button>
        </form>
      </Container>
    </>
  )
}

export default ProfileSetupPage;