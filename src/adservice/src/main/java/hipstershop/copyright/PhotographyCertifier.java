package hipstershop.copyright;

import hipstershop.Demo;

import java.util.ArrayList;
import java.util.List;

class PhotographyCertifier implements Certifier {

    private final StockPhotos stockPhotos = new StockPhotos();

    @Override
    public List<Demo.Ad> certify(List<Demo.Ad> ads) {
        List<Demo.Ad> result = new ArrayList<>(ads.size());
        for (Demo.Ad ad : ads) {
            if(stockPhotos.isCopyright(ad)){
               result.add(ad);
            }
        }
        return result;
    }
}
