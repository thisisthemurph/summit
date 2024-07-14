type PageHeaderProps = {
  title: string | null | undefined;
  subtitle?: string;
}

function PageHeader({title, subtitle}: PageHeaderProps) {
  return (
    <section className="mb-6 shadow-lg p-4">
        <h2 className="text-2xl font-semibold">{title}</h2>
        {subtitle && (
          <p className="text-sm text-slate-800">{subtitle}</p>
        )}
    </section>
  )
}

export default PageHeader;
