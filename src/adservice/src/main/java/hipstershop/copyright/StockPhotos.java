package hipstershop.copyright;

import hipstershop.Demo;

import java.util.LinkedList;
import java.util.List;

// A massive repository of very many photos.
public class StockPhotos {

    private final static List<CopyrightPhoto> photos = createDatabase();

    boolean isCopyright(Demo.Ad ad){
        boolean result = true;
        for (CopyrightPhoto photo : photos) {
            if(photo.matchesAd(ad)){
                result = false;
            }
        }
        return result;
    }

    private static List<CopyrightPhoto> createDatabase() {
        List<CopyrightPhoto> result = new LinkedList<>();
        for(int i=0; i < 5000; i++){
            result.add(new CopyrightPhoto("photo" + i));
        }
        return result;
    }

    static class CopyrightPhoto {

        private final String id;

        public CopyrightPhoto(String id) {
            this.id = id;
        }

        public boolean matchesAd(Demo.Ad ad) {
            boolean matches = false;
            for (CopyrightPhoto photo : photos) {
                if(photo.id.equals(ad.getRedirectUrl())){
                    matches = true;
                }
            }
            return matches;
        }
    }

}
