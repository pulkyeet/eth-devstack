interface Props {
  message: string;
}

export default function ErrorMessage({ message }: Props) {
  return (
    <div className="rounded-lg border border-red-200 bg-red-50 p-4 text-red-800 dark:border-red-900 dark:bg-red-950 dark:text-red-200">
      <p className="font-medium">Error</p>
      <p className="text-sm">{message}</p>
    </div>
  );
}