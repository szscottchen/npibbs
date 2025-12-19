export async function useUploadImage(file: File): Promise<any> {
    const formData = new FormData();
    formData.append("file", file, file.name);
    
    return useHttp("/api/upload", {
        method: "POST",
        body: formData,
    })
}