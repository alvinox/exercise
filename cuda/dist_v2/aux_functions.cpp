#include <math.h>

#include "aux_functions.hh"

float scale(int i, int n) {
  return ((float)i/(n-1));
}

float distance(float x1, float x2) {
  return sqrt((x2-x1)*(x2-x1));
}

void distanceArray(float* out, float* in, float ref, int n) {
  for (int i = 0; i < n; ++i) {
    out[i] = distance(ref, in[i]);
  }
}
