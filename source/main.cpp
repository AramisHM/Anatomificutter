/*	Aramis' Anatomificutter Tool
copyright (c) 2019 by Aramis Hornung Moraes
This file is part of the ahmBitmap project
read ahmbmp.h for copying conditions.
*/

#define AHMBMP_VERSION "0.1.0"

#include <stdio.h>
#include "ahmbmp.h"
#include "math.h"
#include "string.h"
#include "time.h"

/* debug stuff for MSVC */
#ifdef _WIN32_debug
#ifdef _MSC_VER
#define _CRTDBG_MAP_ALLOC
#include <crtdbg.h>
#include <stdlib.h>
#endif
#endif
#ifdef _WIN32
#include "windows.h"
#endif

void main_process_image(unsigned int z_index, unsigned int y_start,
                        unsigned int y_end) {
    ahm_bitmap *myBmp;
    ahm_bitmap *trgt_img;
    char *input_file_name;
    char file_output_name[4096];
    int imHei, imWid;

    if (y_start > y_end) {
        printf("Start index must be smaller than the end index\n");
    }

    char *filename;
    filename = (char *)calloc(30, sizeof(char));
    imHei = y_end - y_start + 1;

    /* Open first image just to get width */
    char *itoa_char_array;
    sprintf(filename, "%d.bmp\0", y_start);
    myBmp = create_bmp_from_file(filename);
    imWid = myBmp->Width;
    destroy_ahmBitmap(myBmp);

    trgt_img = create_ahmBitmap(imWid, imHei);

    /* Iterate each line and recreate the sagital image*/
    for (int i = y_start; i <= y_end; ++i) {
        char *itoa_char_array;
        sprintf(filename, "%d.bmp", i);

        /* Open axial image */
        myBmp = create_bmp_from_file(filename);
        if (!myBmp) {
            destroy_ahmBitmap(myBmp);
            printf(
                "Error: File is not supported. Only 24bit bitmaps can be "
                "processed.\n");
            free(filename);
            return;
        }

        for (int w = 0; w < myBmp->Width; ++w) {
            ahm_pixel p = get_pixel(myBmp, w, z_index);
            int topDown = ((imHei - 1) - (i - 1));
            set_pixel(trgt_img, w, topDown, p.r_, p.g_, p.b_);
        }
        destroy_ahmBitmap(myBmp);
    }
    save_bmp(trgt_img, "out.bmp");
    free(filename);
    destroy_ahmBitmap(trgt_img);
    destroy_ahmBitmap(myBmp);
}

int main(int argc, char *argv[]) {
#ifdef _WIN32_debug
#ifdef _MSC_VER
#ifdef WIN32
    _CrtSetDbgFlag(_CRTDBG_ALLOC_MEM_DF | _CRTDBG_LEAK_CHECK_DF);
#endif
#endif
#endif
    if (argc > 0) {
        printf("do stuff");
        main_process_image(0, 1, 3);
    } else {
        printf("Usage: <usage instructions>\n");
        fflush(stdin);
        getchar();
    }
#ifdef _WIN32_debug
#ifdef _MSC_VER
#ifdef WIN32
    _CrtSetReportMode(_CRT_ERROR, _CRTDBG_MODE_DEBUG);
#endif
#endif
#endif
    return 0;
}
/*
#if defined(_WIN32)
int APIENTRY _WinMain(HINSTANCE hInstance, HINSTANCE hPevInstance, LPSTR
lpCmdLine, int nCmdShow)
{
        return main( __argc, __argv);
}
#endif*/
