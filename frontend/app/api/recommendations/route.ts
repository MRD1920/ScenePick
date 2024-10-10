import { NextResponse } from "next/server";

export async function GET(req: Request) {
  const { searchParams } = new URL(req.url);
  const id = searchParams.get("id");

  const recommendationsResponse = {
    page: 1,
    results: [
      {
        adult: false,
        backdrop_path: "/aFTYFqrWp4RS46Twm87l5e0ItYb.jpg",
        genre_ids: [28, 12],
        id: 2,
        media_type: "movie",
        original_language: "en",
        original_title: "Avatar II",
        overview: "The sequel to Avatar.",
        popularity: 99.0,
        poster_path: "/path/to/poster.jpg",
        release_date: "2025-12-18",
        title: "Avatar II",
        video: false,
        vote_average: 8.0,
        vote_count: 5000,
      },
      // Add more recommendations as needed
    ],
    total_pages: 1,
    total_results: 1,
  };

  return NextResponse.json(recommendationsResponse);
}
